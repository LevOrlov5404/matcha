package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/l-orlov/matcha/internal/service"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	ctxUserID   = "userID"
	ctxLogEntry = "log-entry"
)

var ErrNotValidAuthorizationHeader = errors.New("not valid Authorization header")

func (s *Server) InitMiddleware(c *gin.Context) {
	requestID := uuid.New().String()
	logEntry := logrus.NewEntry(s.log).WithField("request-id", requestID)
	c.Set(ctxLogEntry, logEntry)
}

func (s *Server) UserAuthorizationMiddleware(c *gin.Context) {
	err := s.validateTokenCookieAndRefreshIfNeeded(c)
	if err == nil {
		return
	}
	getLogEntry(c).Debug(err)

	if err := s.validateTokenHeader(c); err != nil {
		s.newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}
}

// validateTokenCookieAndRefreshIfNeeded gets accessToken from cookie and validate it.
// on success it puts accessToken data to ctx and returns nil.
// else it tries to refresh session by refresh token from cookie:
// - on success puts accessToken data to ctx and returns nil
// - on failure returns error.
func (s *Server) validateTokenCookieAndRefreshIfNeeded(c *gin.Context) error {
	accessToken, err := s.Cookie(c, accessTokenCookieName)
	if err != nil {
		getLogEntry(c).Debug(err)
		return s.refreshSessionByRefreshTokenCookie(c)
	}

	accessTokenClaims, err := s.svc.UserAuthorization.ValidateAccessToken(accessToken)
	if err != nil {
		if !strings.Contains(err.Error(), "token is expired by") &&
			!errors.Is(err, service.ErrNotActiveAccessToken) {
			return err
		}

		return s.refreshSessionByRefreshTokenCookie(c)
	}

	return setUserIDForContext(c, accessTokenClaims.Subject)
}

func (s *Server) refreshSessionByRefreshTokenCookie(c *gin.Context) error {
	refreshToken, err := s.Cookie(c, refreshTokenCookieName)
	if err != nil {
		return err
	}

	newAccessToken, newRefreshToken, err := s.svc.UserAuthorization.RefreshSession(refreshToken)
	if err != nil {
		return err
	}

	accessTokenClaims, err := s.svc.UserAuthorization.GetAccessTokenClaims(newAccessToken)
	if err != nil {
		return err
	}

	s.setTokensCookies(c, newAccessToken, newRefreshToken)

	return setUserIDForContext(c, accessTokenClaims.Subject)
}

// validateTokenHeader gets accessToken from header and validate it.
// on success it puts accessToken data to ctx and returns nil. else it returns error.
func (s *Server) validateTokenHeader(c *gin.Context) error {
	header := c.GetHeader("Authorization")
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return ErrNotValidAuthorizationHeader
	}

	accessToken := headerParts[1]
	accessTokenClaims, err := s.svc.UserAuthorization.ValidateAccessToken(accessToken)
	if err != nil {
		return err
	}

	return setUserIDForContext(c, accessTokenClaims.Subject)
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

func setUserIDForContext(c *gin.Context, userIDStr string) error {
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return err
	}

	c.Set(ctxUserID, userID)

	return nil
}
