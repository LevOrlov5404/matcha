package server

import (
	"errors"
	ierrors "github.com/LevOrlov5404/matcha/internal/errors"
	"net/http"
	"strconv"

	"github.com/LevOrlov5404/matcha/internal/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateUser(c *gin.Context) {
	var user models.UserToCreate
	if err := c.BindJSON(&user); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, ierrors.NewBusiness(err, ""))
		return
	}

	id, err := s.services.User.CreateUser(c, user)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (s *Server) SignIn(c *gin.Context) {
	var user models.UserToSignIn
	if err := c.BindJSON(&user); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, ierrors.NewBusiness(err, ""))
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

func (s *Server) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		s.newErrorResponse(
			c, http.StatusBadRequest, ierrors.NewBusiness(errors.New("invalid id param"), ""),
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
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, ierrors.NewBusiness(err, ""))
		return
	}

	if err := s.services.User.UpdateUser(c, user); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) GetAllUsers(c *gin.Context) {
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		s.newErrorResponse(
			c, http.StatusBadRequest, ierrors.NewBusiness(errors.New("invalid id param"), ""),
		)
		return
	}

	if err := s.services.User.DeleteUser(c, id); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
