package server

import (
	"net/http"
	"strconv"

	"github.com/LevOrlov5404/matcha/internal/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateUser(c *gin.Context) {
	var user models.UserToCreate
	if err := c.BindJSON(&user); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := s.services.User.CreateUser(c, user)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (s *Server) SignIn(c *gin.Context) {
	var user models.UserToSignIn
	if err := c.BindJSON(&user); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := s.services.User.GenerateToken(c, user.Email, user.Password)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (s *Server) GetUserByEmailPassword(c *gin.Context) {
	var inputUser models.UserToSignIn
	if err := c.BindJSON(&inputUser); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := s.services.User.GetUserByEmailPassword(c, inputUser.Email, inputUser.Password)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	user, err := s.services.User.GetUserByID(c, id)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var user models.UserToCreate
	if err := c.BindJSON(&user); err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.services.User.UpdateUser(c, id, user); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) GetAllUsers(c *gin.Context) {
	users, err := s.services.User.GetAllUsers(c)
	if err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, users)
}

func (s *Server) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		s.newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := s.services.User.DeleteUser(c, id); err != nil {
		s.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
