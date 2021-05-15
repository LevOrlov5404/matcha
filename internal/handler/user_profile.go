package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ierrors "github.com/l-orlov/matcha/internal/errors"
	"github.com/l-orlov/matcha/internal/models"
)

func (h *Handler) GetUserProfileByID(c *gin.Context) {
	setHandlerNameToLogEntry(c, "GetUserProfileByID")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		h.newErrorResponse(
			c, http.StatusBadRequest, ierrors.NewBusiness(ErrNotValidIDParameter, ""),
		)
		return
	}

	user, err := h.svc.UserProfile.GetUserProfileByID(c, id)
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

func (h *Handler) UpdateUserProfile(c *gin.Context) {
	setHandlerNameToLogEntry(c, "UpdateUserProfile")

	var user models.UserProfile
	if err := c.BindJSON(&user); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.UserProfile.UpdateUserProfile(c, user); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UploadUserAvatar(c *gin.Context) {
	setHandlerNameToLogEntry(c, "UploadUserAvatar")

	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if err = h.svc.UserProfile.UploadUserAvatar(c, userID, file); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) DeleteUserAvatar(c *gin.Context) {
	setHandlerNameToLogEntry(c, "DeleteUserAvatar")

	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if err := h.svc.UserProfile.DeleteUserAvatar(c, userID); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) UploadUserPicture(c *gin.Context) {
	setHandlerNameToLogEntry(c, "UploadUserPicture")

	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if err = h.svc.UserProfile.UploadUserPicture(c, userID, file); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetUserPictures(c *gin.Context) {
	setHandlerNameToLogEntry(c, "GetUserPictures")

	userID, err := getUserIDFromContext(c)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	users, err := h.svc.UserProfile.GetUserPicturesByUserID(c, userID)
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

func (h *Handler) DeleteUserPicture(c *gin.Context) {
	setHandlerNameToLogEntry(c, "DeleteUserPicture")

	pictureUUIDStr := c.Query("uuid")
	if pictureUUIDStr == "" {
		h.newErrorResponse(c, http.StatusBadRequest, ErrNotValidUUIDParameter)
		return
	}

	pictureUUID, err := uuid.Parse(pictureUUIDStr)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.UserProfile.DeleteUserPicture(c, pictureUUID); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
