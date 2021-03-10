package server

import (
	"net/http"

	ierrors "github.com/LevOrlov5404/matcha/internal/errors"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (s *Server) newErrorResponse(c *gin.Context, statusCode int, err error) {
	if customErr, ok := err.(*ierrors.Error); ok {
		s.handleCustomError(c, customErr)
		return
	}

	s.log.Error(err)

	errResp := &errorResponse{
		Message: err.Error(),
	}
	if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		errResp.Detail = ierrors.DetailBusiness
	} else {
		errResp.Detail = ierrors.DetailServer
	}

	c.AbortWithStatusJSON(statusCode, errResp)
}

func (s *Server) handleCustomError(c *gin.Context, err *ierrors.Error) {
	var statusCode int

	if err.Level == ierrors.Business {
		s.log.Debug(err)
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
