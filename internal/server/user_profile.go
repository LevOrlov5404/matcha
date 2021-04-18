package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	ierrors "github.com/l-orlov/matcha/internal/errors"
	"github.com/l-orlov/matcha/internal/models"
)

func (s *Server) GetUserProfileByID(c *gin.Context) {
	setHandlerNameToLogEntry(c, "GetUserProfileByID")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		s.newErrorResponse(
			c, http.StatusBadRequest, ierrors.NewBusiness(ErrNotValidIDParameter, ""),
		)
		return
	}

	user, err := s.svc.UserProfile.GetUserProfileByID(c, id)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) UpdateUserProfile(c *gin.Context) {
	setHandlerNameToLogEntry(c, "UpdateUserProfile")

	var user models.UserProfile
	if err := c.BindJSON(&user); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := s.svc.UserProfile.UpdateUserProfile(c, user); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) UploadUserAvatar(c *gin.Context) {
	setHandlerNameToLogEntry(c, "UploadUserAvatar")

	userID, ok := getUserIDFromContext(c)
	if !ok {
		s.newErrorResponse(c, http.StatusInternalServerError, ErrNotValidIDParameter)
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if err = s.svc.UserProfile.UploadUserAvatar(c, userID, file); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) DeleteUserAvatar(c *gin.Context) {
	setHandlerNameToLogEntry(c, "DeleteUserAvatar")

	userID, ok := getUserIDFromContext(c)
	if !ok {
		s.newErrorResponse(c, http.StatusInternalServerError, ErrNotValidIDParameter)
		return
	}

	if err := s.svc.UserProfile.DeleteUserAvatar(c, userID); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) UploadUserPicture(c *gin.Context) {
	setHandlerNameToLogEntry(c, "UploadUserPicture")

	userID, ok := getUserIDFromContext(c)
	if !ok {
		s.newErrorResponse(c, http.StatusInternalServerError, ErrNotValidIDParameter)
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if err = s.svc.UserProfile.UploadUserPicture(c, userID, file); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) GetUserPictures(c *gin.Context) {
	setHandlerNameToLogEntry(c, "GetUserPictures")

	userID, ok := getUserIDFromContext(c)
	if !ok {
		s.newErrorResponse(c, http.StatusInternalServerError, ErrNotValidIDParameter)
		return
	}

	users, err := s.svc.UserProfile.GetUserPicturesByUserID(c, userID)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if users == nil {
		c.JSON(http.StatusOK, []struct{}{})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (s *Server) DeleteUserPicture(c *gin.Context) {
	setHandlerNameToLogEntry(c, "DeleteUserPicture")

	pictureUUIDStr := c.Query("uuid")
	if pictureUUIDStr == "" {
		s.newErrorResponse(c, http.StatusBadRequest, ErrNotValidUUIDParameter)
		return
	}

	pictureUUID, err := uuid.Parse(pictureUUIDStr)
	if err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := s.svc.UserProfile.DeleteUserPicture(c, pictureUUID); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func getUserIDFromContext(c *gin.Context) (uint64, bool) {
	userIDParam, ok := c.Get(ctxUserID)
	if !ok {
		return 0, false
	}

	userID, ok := userIDParam.(uint64)
	if !ok {
		return 0, false
	}

	return userID, true
}
