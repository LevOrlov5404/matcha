package server

import (
	"net/http"

	iErrs "github.com/LevOrlov5404/matcha/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (s *Server) newErrorResponse(c *gin.Context, statusCode int, err error) {
	logEntry := getLogEntry(c)

	if customErr, ok := err.(*iErrs.Error); ok {
		handleCustomError(c, customErr, logEntry)
		return
	}

	logEntry.Error(err)

	errResp := &errorResponse{
		Message: err.Error(),
	}
	if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		errResp.Detail = iErrs.DetailBusiness
	} else {
		errResp.Detail = iErrs.DetailServer
	}

	c.AbortWithStatusJSON(statusCode, errResp)
}

func handleCustomError(c *gin.Context, err *iErrs.Error, logEntry *logrus.Entry) {
	var statusCode int

	if err.Level == iErrs.Business {
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
