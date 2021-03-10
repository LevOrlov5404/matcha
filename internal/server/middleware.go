package server

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	headerAuth = "Authorization"
	ctxUser    = "userID"

	keyLogEntry = "log-entry"
)

func (s *Server) InitMiddleware(c *gin.Context) {
	start := time.Now()

	requestID := uuid.New().String()
	logEntry := logrus.NewEntry(s.log).WithField("request-id", requestID)
	logEntry.Debugf("Start request at %s", start.Format(time.RFC1123Z))

	c.Set(keyLogEntry, logEntry)

	c.Next()

	logEntry = getLogEntry(c)
	logEntry.WithField("elapsed", time.Since(start).String()).
		Debugf("Complete request at %s", time.Now().Format(time.RFC1123Z))
}

func (s *Server) UserIdentityMiddleware(c *gin.Context) {
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

func setHandlerNameToLogEntry(c *gin.Context, handlerName string) {
	logEntryValue, _ := c.Get(keyLogEntry)

	logEntry := logEntryValue.(*logrus.Entry).WithField("method", handlerName)
	c.Set(keyLogEntry, logEntry)
}

func getLogEntry(c *gin.Context) *logrus.Entry {
	logEntryValue, _ := c.Get(keyLogEntry)

	return logEntryValue.(*logrus.Entry)
}
