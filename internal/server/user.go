package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	iErrs "github.com/l-orlov/matcha/internal/errors"
	"github.com/l-orlov/matcha/internal/models"
	"github.com/pkg/errors"
)

func (s *Server) CreateUser(c *gin.Context) {
	setHandlerNameToLogEntry(c, "CreateUser")

	var user models.UserToCreate
	if err := c.BindJSON(&user); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := s.svc.User.CreateUser(c, user)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	emailConfirmToken, err := s.svc.Verification.CreateEmailConfirmToken(id)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	// send token by email
	s.svc.Mailer.SendEmailConfirm(user.Email, emailConfirmToken)

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (s *Server) GetUserByID(c *gin.Context) {
	setHandlerNameToLogEntry(c, "GetUserByID")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		s.newErrorResponse(
			c, http.StatusBadRequest, iErrs.NewBusiness(errors.New("invalid id parameter"), ""),
		)
		return
	}

	user, err := s.svc.User.GetUserByID(c, id)
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

func (s *Server) UpdateUser(c *gin.Context) {
	setHandlerNameToLogEntry(c, "UpdateUser")

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := s.svc.User.UpdateUser(c, user); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) SetUserPassword(c *gin.Context) {
	setHandlerNameToLogEntry(c, "SetPassword")

	var user models.UserPassword
	if err := c.BindJSON(&user); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := s.svc.User.SetUserPassword(c, user.ID, user.Password); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) ChangeUserPassword(c *gin.Context) {
	setHandlerNameToLogEntry(c, "ChangePassword")

	var user models.UserPasswordToChange
	if err := c.BindJSON(&user); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := s.svc.User.ChangeUserPassword(c, user.ID, user.OldPassword, user.NewPassword); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) GetAllUsers(c *gin.Context) {
	setHandlerNameToLogEntry(c, "GetAllUsers")

	users, err := s.svc.User.GetAllUsers(c)
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

func (s *Server) DeleteUser(c *gin.Context) {
	setHandlerNameToLogEntry(c, "DeleteUser")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		s.newErrorResponse(
			c, http.StatusBadRequest, iErrs.NewBusiness(errors.New("invalid id parameter"), ""),
		)
		return
	}

	if err := s.svc.User.DeleteUser(c, id); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) GetUserProfileByID(c *gin.Context) {
	setHandlerNameToLogEntry(c, "GetUserProfileByID")

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		s.newErrorResponse(
			c, http.StatusBadRequest, iErrs.NewBusiness(errors.New("invalid id parameter"), ""),
		)
		return
	}

	user, err := s.svc.User.GetUserProfileByID(c, id)
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

	if err := s.svc.User.UpdateUserProfile(c, user); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
