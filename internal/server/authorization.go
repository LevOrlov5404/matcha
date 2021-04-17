package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	ierrors "github.com/l-orlov/matcha/internal/errors"
	"github.com/l-orlov/matcha/internal/models"
	"github.com/pkg/errors"
)

const (
	accessTokenCookieName  = "access_token"
	refreshTokenCookieName = "refresh_token"
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

	accessToken, refreshToken, err := s.svc.CreateSession(strconv.FormatUint(userID, 10))
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	s.setTokensCookies(c, accessToken, refreshToken)
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

	if _, err := s.svc.UserAuthorization.ValidateAccessToken(req.AccessToken); err != nil {
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

	accessToken, refreshToken, err := s.svc.RefreshSession(req.RefreshToken)
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

	accessToken, err := s.Cookie(c, accessTokenCookieName)
	if err != nil {
		getLogEntry(c).Debug(err)
	}

	if accessToken == "" {
		// try to get accessToken from request
		var req models.LogoutRequest
		if err := c.BindJSON(&req); err != nil {
			s.newErrorResponse(c, http.StatusBadRequest, err)
			return
		}

		accessToken = req.AccessToken
	}

	if err = s.svc.RevokeSession(accessToken); err != nil {
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
			c, http.StatusBadRequest, ierrors.NewBusiness(ErrEmptyEmailParameter, ""),
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
			c, http.StatusBadRequest, ierrors.NewBusiness(ErrUserNotFound, ""),
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

func (s *Server) setTokensCookies(c *gin.Context, accessToken, refreshToken string) {
	if encodedAccessToken, err := s.options.SecureCookie.Encode(accessTokenCookieName, accessToken); err == nil {
		c.SetCookie(
			accessTokenCookieName, encodedAccessToken, s.options.AccessTokenCookieMaxAge,
			"/", s.cfg.Cookie.Domain, false, true,
		)
	} else {
		getLogEntry(c).Error(err)
	}

	if encodedRefreshToken, err := s.options.SecureCookie.Encode(refreshTokenCookieName, refreshToken); err == nil {
		c.SetCookie(
			refreshTokenCookieName, encodedRefreshToken, s.options.RefreshTokenCookieMaxAge,
			"/", s.cfg.Cookie.Domain, false, true,
		)
	} else {
		getLogEntry(c).Error(err)
	}
}

func (s *Server) Cookie(c *gin.Context, name string) (string, error) {
	value, err := c.Cookie(name)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get cookie %s", name)
	}

	var decodedValue string
	if err := s.options.SecureCookie.Decode(name, value, &decodedValue); err != nil {
		return "", err
	}

	return decodedValue, nil
}
