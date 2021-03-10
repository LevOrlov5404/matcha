package server

import (
	"net/http"

	ierrors "github.com/LevOrlov5404/matcha/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (s *Server) newErrorResponse(c *gin.Context, statusCode int, err error) {
	logEntry := getLogEntry(c)

	if customErr, ok := err.(*ierrors.Error); ok {
		handleCustomError(c, customErr, logEntry)
		return
	}

	logEntry.Error(err)

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

func handleCustomError(c *gin.Context, err *ierrors.Error, logEntry *logrus.Entry) {
	var statusCode int

	if err.Level == ierrors.Business {
		logEntry.Debug(err)
		statusCode = http.StatusBadRequest
	} else {
		logEntry.Error(err)
		statusCode = http.StatusInternalServerError
	}

	c.AbortWithStatusJSON(statusCode, &errorResponse{
		Message: err.Error(),
		Detail:  err.Detail,
	})
}
