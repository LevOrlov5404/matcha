package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		s.newErrorResponse(
			c, http.StatusBadRequest, ierrors.NewBusiness(ErrNotValidIDParameter, ""),
		)
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

	err = s.svc.UserProfile.UploadUserAvatar(c, userID, file)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) UploadUserPicture(c *gin.Context) {
	setHandlerNameToLogEntry(c, "UploadUserPicture")

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		s.newErrorResponse(
			c, http.StatusBadRequest, ierrors.NewBusiness(ErrNotValidIDParameter, ""),
		)
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

	err = s.svc.UserProfile.UploadUserPicture(c, userID, file)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
