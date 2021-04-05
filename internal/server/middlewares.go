package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/l-orlov/matcha/internal/service"
	"github.com/sirupsen/logrus"
)

const (
	ctxUserID   = "userID"
	ctxLogEntry = "log-entry"
)

func (s *Server) InitMiddleware(c *gin.Context) {
	requestID := uuid.New().String()
	logEntry := logrus.NewEntry(s.log).WithField("request-id", requestID)
	c.Set(ctxLogEntry, logEntry)
}

func (s *Server) UserAuthorizationMiddleware(c *gin.Context) {
	accessToken, err := s.Cookie(c, accessTokenCookieName)
	if err != nil {
		getLogEntry(c).Debug(err)
		s.refreshSessionByRefreshTokenCookie(c)
		return
	}

	accessTokenClaims, err := s.svc.UserAuthorization.ValidateAccessToken(accessToken)
	if err != nil {
		if !errors.Is(err, service.ErrNotActiveAccessToken) {
			s.newErrorResponse(c, http.StatusUnauthorized, err)
			return
		}

		s.refreshSessionByRefreshTokenCookie(c)
		return
	}

	c.Set(ctxUserID, accessTokenClaims.Subject)
}

func (s *Server) refreshSessionByRefreshTokenCookie(c *gin.Context) {
	refreshToken, err := s.Cookie(c, refreshTokenCookieName)
	if err != nil {
		s.newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	newAccessToken, newRefreshToken, err := s.svc.UserAuthorization.RefreshSession(refreshToken)
	if err != nil {
		s.newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	accessTokenClaims, err := s.svc.UserAuthorization.GetAccessTokenClaims(newAccessToken)
	if err != nil {
		s.newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	s.setTokensCookies(c, newAccessToken, newRefreshToken)
	c.Set(ctxUserID, accessTokenClaims.Subject)
}

func setHandlerNameToLogEntry(c *gin.Context, handlerName string) {
	logEntryValue, _ := c.Get(ctxLogEntry)

	logEntry := logEntryValue.(*logrus.Entry).WithField("method", handlerName)
	c.Set(ctxLogEntry, logEntry)
}

func getLogEntry(c *gin.Context) *logrus.Entry {
	logEntryValue, _ := c.Get(ctxLogEntry)

	return logEntryValue.(*logrus.Entry)
}
