package server

import (
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

func (s *Server) newErrorResponse(c *gin.Context, statusCode int, message string) {
	s.log.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{
		Message: message,
	})
}
