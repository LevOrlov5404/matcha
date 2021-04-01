package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	iErrs "github.com/l-orlov/matcha/internal/errors"
	"github.com/l-orlov/matcha/internal/models"
	"github.com/pkg/errors"
)

func (s *Server) SignIn(c *gin.Context) {
	setHandlerNameToLogEntry(c, "SignIn")

	var user models.UserToSignIn
	if err := c.BindJSON(&user); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	userID, err := s.svc.AuthenticateUserByUsername(c, user.Username, user.Password, user.Fingerprint)
	if err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	accessToken, refreshToken, err := s.svc.CreateSession(strconv.FormatUint(userID, 10), user.Fingerprint)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (s *Server) ValidateAccessToken(c *gin.Context) {
	setHandlerNameToLogEntry(c, "ValidateAccessToken")

	var req models.ValidateAccessTokenRequest
	if err := c.BindJSON(&req); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	_, err := s.svc.UserAuthorization.ValidateAccessToken(req.AccessToken)
	if err != nil {
		s.newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) RefreshSession(c *gin.Context) {
	setHandlerNameToLogEntry(c, "RefreshSession")

	var req models.RefreshSessionRequest
	if err := c.BindJSON(&req); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	accessToken, refreshToken, err := s.svc.RefreshSession(req.RefreshToken, req.Fingerprint)
	if err != nil {
		s.newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (s *Server) Logout(c *gin.Context) {
	setHandlerNameToLogEntry(c, "Logout")

	var req models.LogoutRequest
	if err := c.BindJSON(&req); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err := s.svc.RevokeSession(req.AccessToken)
	if err != nil {
		s.newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) ResetPassword(c *gin.Context) {
	setHandlerNameToLogEntry(c, "ResetPassword")

	email, ok := c.GetQuery("email")
	if !ok || email == "" {
		s.newErrorResponse(
			c, http.StatusBadRequest, iErrs.NewBusiness(errors.New("empty email parameter"), ""),
		)
		return
	}

	user, err := s.svc.User.GetUserByEmail(c, email)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		s.newErrorResponse(
			c, http.StatusBadRequest, iErrs.NewBusiness(errors.New("there is no user with this email"), ""),
		)
		return
	}

	passwordResetConfirmToken, err := s.svc.Verification.CreatePasswordResetConfirmToken(user.ID)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	// send token by email
	s.svc.Mailer.SendResetPasswordConfirm(user.Email, passwordResetConfirmToken)

	c.Status(http.StatusOK)
}
