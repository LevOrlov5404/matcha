package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	iErrs "github.com/l-orlov/matcha/internal/errors"
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

	userID, err := s.svc.Verification.VerifyEmailConfirmToken(token)
	if err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := s.svc.User.ConfirmEmail(c, userID); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) ConfirmPasswordReset(c *gin.Context) {
	setHandlerNameToLogEntry(c, "ConfirmPasswordReset")

	token, ok := c.GetQuery("token")
	if !ok || token == "" {
		s.newErrorResponse(
			c, http.StatusBadRequest, iErrs.NewBusiness(errors.New("empty token parameter"), ""),
		)
		return
	}

	userID, err := s.svc.Verification.VerifyPasswordResetConfirmToken(token)
	if err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": userID,
	})
}
