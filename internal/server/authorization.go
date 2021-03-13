package server

import (
	"net/http"

	iErrs "github.com/LevOrlov5404/matcha/internal/errors"
	"github.com/LevOrlov5404/matcha/internal/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) SignIn(c *gin.Context) {
	setHandlerNameToLogEntry(c, "SignIn")

	var user models.UserToSignIn
	if err := c.BindJSON(&user); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, iErrs.NewBusiness(err, ""))
		return
	}

	token, err := s.services.User.GenerateToken(c, user.Username, user.Password)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (s *Server) ResetPassword(c *gin.Context) {
	setHandlerNameToLogEntry(c, "ResetPassword")

	var userEmail models.UserEmail
	if err := c.BindJSON(&userEmail); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, iErrs.NewBusiness(err, ""))
		return
	}

	user, err := s.services.User.GetUserByEmail(c, userEmail.Email)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	resetPasswordConfirmToken, err := s.services.Verification.CreateResetPasswordConfirmToken(user.ID)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	// send token by email
	if err := s.services.Mailer.SendResetPasswordConfirm(user.Email, resetPasswordConfirmToken); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
