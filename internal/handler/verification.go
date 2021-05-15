package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ierrors "github.com/l-orlov/matcha/internal/errors"
)

func (h *Handler) ConfirmEmail(c *gin.Context) {
	setHandlerNameToLogEntry(c, "ConfirmEmail")

	token, ok := c.GetQuery("token")
	if !ok || token == "" {
		h.newErrorResponse(
			c, http.StatusBadRequest, ierrors.NewBusiness(ErrEmptyTokenParameter, ""),
		)
		return
	}

	userID, err := h.svc.Verification.VerifyEmailConfirmToken(token)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.User.ConfirmEmail(c, userID); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) ConfirmPasswordReset(c *gin.Context) {
	setHandlerNameToLogEntry(c, "ConfirmPasswordReset")

	token, ok := c.GetQuery("token")
	if !ok || token == "" {
		h.newErrorResponse(
			c, http.StatusBadRequest, ierrors.NewBusiness(ErrEmptyTokenParameter, ""),
		)
		return
	}

	userID, err := h.svc.Verification.VerifyPasswordResetConfirmToken(token)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": userID,
	})
}
