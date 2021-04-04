package server

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/l-orlov/matcha/internal/service"
	"github.com/sirupsen/logrus"
)

const (
	headerAuth = "Authorization"
	ctxUser    = "userID"

	keyLogEntry = "log-entry"
)

func (s *Server) InitMiddleware(c *gin.Context) {
	requestID := uuid.New().String()
	logEntry := logrus.NewEntry(s.log).WithField("request-id", requestID)
	c.Set(keyLogEntry, logEntry)
}

func (s *Server) UserAuthorizationMiddleware(c *gin.Context) {
	accessToken, err := c.Cookie(accessTokenCookieName)
	if err != nil {
		getLogEntry(c).Debug(err)
	}

	if accessToken == "" {
		// try to get accessToken from header
		header := c.GetHeader(headerAuth)
		if header != "" {
			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				s.newErrorResponse(c, http.StatusUnauthorized, errors.New("invalid auth header"))
				return
			}

			accessToken = headerParts[1]
		}
	}

	accessTokenClaims, err := s.svc.UserAuthorization.ValidateAccessToken(accessToken)
	if err != nil {
		if !errors.Is(err, service.ErrNotActiveAccessToken) {
			s.newErrorResponse(c, http.StatusUnauthorized, err)
			return
		}

		refreshToken, err := c.Cookie(refreshTokenCookieName)
		if err != nil {
			s.newErrorResponse(c, http.StatusUnauthorized, err)
			return
		}

		newAccessToken, newRefreshToken, err := s.svc.UserAuthorization.RefreshSession(refreshToken)
		if err != nil {
			s.newErrorResponse(c, http.StatusUnauthorized, err)
			return
		}

		accessTokenClaims, err = s.svc.UserAuthorization.GetAccessTokenClaims(newAccessToken)
		if err != nil {
			s.newErrorResponse(c, http.StatusUnauthorized, err)
			return
		}

		s.setTokensCookies(c, newAccessToken, newRefreshToken)
	}

	c.Set(ctxUser, accessTokenClaims.Subject)
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
