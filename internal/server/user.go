package server

import (
	"net/http"
	"strconv"

	iErrs "github.com/LevOrlov5404/matcha/internal/errors"
	"github.com/LevOrlov5404/matcha/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (s *Server) CreateUser(c *gin.Context) {
	setHandlerNameToLogEntry(c, "CreateUser")

	var user models.UserToCreate
	if err := c.BindJSON(&user); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, iErrs.NewBusiness(err, ""))
		return
	}

	id, err := s.services.User.CreateUser(c, user)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	// var id uint64 = 0

	emailConfirmToken, err := s.services.Verification.CreateEmailConfirmToken(id)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	// send token by email
	if err := s.services.Mailer.SendEmailConfirm(user.Email, emailConfirmToken); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

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

	user, err := s.services.User.GetUserByID(c, id)
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
		s.newErrorResponse(c, http.StatusBadRequest, iErrs.NewBusiness(err, ""))
		return
	}

	if err := s.services.User.UpdateUser(c, user); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) GetAllUsers(c *gin.Context) {
	setHandlerNameToLogEntry(c, "GetAllUsers")

	users, err := s.services.User.GetAllUsers(c)
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

	if err := s.services.User.DeleteUser(c, id); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
