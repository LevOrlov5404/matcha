package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	ierrors "github.com/l-orlov/matcha/internal/errors"
	"github.com/l-orlov/matcha/internal/models"
)

func (h *Handler) CreateUser(c *gin.Context) {
	setHandlerNameToLogEntry(c, "CreateUser")

	var user models.UserToCreate
	if err := c.BindJSON(&user); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := h.svc.User.CreateUser(c, user)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	emailConfirmToken, err := h.svc.Verification.CreateEmailConfirmToken(id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	// send token by email
	h.svc.Mailer.SendEmailConfirm(user.Email, emailConfirmToken)

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) GetUserByID(c *gin.Context) {
	setHandlerNameToLogEntry(c, "GetUserByID")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(
			c, http.StatusBadRequest, ierrors.NewBusiness(ErrNotValidIDParameter, ""),
		)
		return
	}

	user, err := h.svc.User.GetUserByID(c, id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	setHandlerNameToLogEntry(c, "UpdateUser")

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.User.UpdateUser(c, user); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) SetUserPassword(c *gin.Context) {
	setHandlerNameToLogEntry(c, "SetPassword")

	var user models.UserPassword
	if err := c.BindJSON(&user); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	// ToDo: add deleting all sessions

	if err := h.svc.User.SetUserPassword(c, user.ID, user.Password); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) ChangeUserPassword(c *gin.Context) {
	setHandlerNameToLogEntry(c, "ChangePassword")

	var user models.UserPasswordToChange
	if err := c.BindJSON(&user); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.User.ChangeUserPassword(c, user.ID, user.OldPassword, user.NewPassword); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	setHandlerNameToLogEntry(c, "GetAllUsers")

	users, err := h.svc.User.GetAllUsers(c)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if users == nil {
		c.JSON(http.StatusOK, []struct{}{})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	setHandlerNameToLogEntry(c, "DeleteUser")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(
			c, http.StatusBadRequest, ierrors.NewBusiness(ErrNotValidIDParameter, ""),
		)
		return
	}

	if err := h.svc.User.DeleteUser(c, id); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
