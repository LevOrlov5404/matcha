package server

import (
	"net/http"

	customErrs "github.com/LevOrlov5404/matcha/internal/custom-errors"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (s *Server) newErrorResponse(c *gin.Context, statusCode int, err error) {
	if customErr, ok := err.(*customErrs.Error); ok {
		s.handleCustomError(c, customErr)
		return
	}

	s.log.Error(err)

	errResp := &errorResponse{
		Message: err.Error(),
	}
	if statusCode >= 400 && statusCode < 500 {
		errResp.Detail = customErrs.DetailBusiness
	} else {
		errResp.Detail = customErrs.DetailServer
	}

	c.AbortWithStatusJSON(statusCode, errResp)
}

func (s *Server) handleCustomError(c *gin.Context, err *customErrs.Error) {
	var statusCode int

	if err.Level == customErrs.Business {
		statusCode = http.StatusBadRequest
	} else {
		s.log.Error(err)

		statusCode = http.StatusInternalServerError
	}

	c.AbortWithStatusJSON(statusCode, &errorResponse{
		Message: err.Error(),
		Detail:  err.Detail,
	})
}
