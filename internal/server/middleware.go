package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	headerAuth = "Authorization"
	ctxUser    = "userID"
)

func (s *Server) UserIdentity(c *gin.Context) {
	header := c.GetHeader(headerAuth)
	if header == "" {
		s.newErrorResponse(c, http.StatusUnauthorized, errors.New("empty auth header"))
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		s.newErrorResponse(c, http.StatusUnauthorized, errors.New("invalid auth header"))
		return
	}

	userID, err := s.services.User.ParseToken(headerParts[1])
	if err != nil {
		s.newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.Set(ctxUser, userID)
}

func getUserID(c *gin.Context) (int64, error) {
	id, ok := c.Get(ctxUser)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int64)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
