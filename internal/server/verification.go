package server

import (
	"net/http"

	iErrs "github.com/LevOrlov5404/matcha/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) ConfirmEmail(c *gin.Context) {
	setHandlerNameToLogEntry(c, "ConfirmEmail")

	token, ok := c.GetQuery("token")
	if !ok || token == "" {
		s.newErrorResponse(
			c, http.StatusBadRequest, iErrs.NewBusiness(errors.New("empty token parameter"), ""),
		)
		return
	}

	userID, err := s.services.Verification.VerifyEmailConfirmToken(token)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if err := s.services.User.ConfirmEmail(c, userID); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) ConfirmResetPassword(c *gin.Context) {
	setHandlerNameToLogEntry(c, "ConfirmResetPassword")

	token, ok := c.GetQuery("token")
	if !ok || token == "" {
		s.newErrorResponse(
			c, http.StatusBadRequest, iErrs.NewBusiness(errors.New("empty token parameter"), ""),
		)
		return
	}

	if _, err := s.services.Verification.VerifyResetPasswordConfirmToken(token); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
