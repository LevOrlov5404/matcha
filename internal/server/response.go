package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	iErrs "github.com/l-orlov/matcha/internal/errors"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (s *Server) newErrorResponse(c *gin.Context, statusCode int, err error) {
	logEntry := getLogEntry(c)

	if customErr, ok := err.(*iErrs.Error); ok {
		handleCustomError(c, logEntry, customErr)
		return
	}

	handleDefaultError(c, logEntry, err, statusCode)
}

func handleCustomError(c *gin.Context, logEntry *logrus.Entry, err *iErrs.Error) {
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

func handleDefaultError(c *gin.Context, logEntry *logrus.Entry, err error, statusCode int) {
	errResp := &errorResponse{
		Message: err.Error(),
	}
	if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		logEntry.Debug(err)
		errResp.Detail = iErrs.DetailBusiness
	} else {
		logEntry.Error(err)
		errResp.Detail = iErrs.DetailServer
	}

	c.AbortWithStatusJSON(statusCode, errResp)
}
