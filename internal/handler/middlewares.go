package handler

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

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func (h *Handler) InitMiddleware(c *gin.Context) {
	requestID := uuid.New().String()
	logEntry := logrus.NewEntry(h.log).WithField("request-id", requestID)
	c.Set(ctxLogEntry, logEntry)
}

func (h *Handler) UserAuthorizationMiddleware(c *gin.Context) {
	if err := h.validateTokenCookieAndRefreshIfNeeded(c); err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.Next()
}

// validateTokenCookieAndRefreshIfNeeded gets accessToken from cookie and validate it.
// on success it puts accessToken data to ctx and returns nil.
// else it tries to refresh session by refresh token from cookie:
// - on success puts accessToken data to ctx and returns nil
// - on failure returns error.
func (h *Handler) validateTokenCookieAndRefreshIfNeeded(c *gin.Context) error {
	accessToken, err := h.Cookie(c, accessTokenCookieName)
	if err != nil {
		h.getLogEntry(c).Debug(err)
		return h.refreshSessionByRefreshTokenCookie(c)
	}

	accessTokenClaims, err := h.svc.UserAuthorization.ValidateAccessToken(accessToken)
	if err != nil {
		if !strings.Contains(err.Error(), "token is expired by") &&
			!errors.Is(err, service.ErrNotActiveAccessToken) {
			return err
		}

		return h.refreshSessionByRefreshTokenCookie(c)
	}

	return h.validateAndSetUserIDForContext(c, accessTokenClaims.Subject)
}

func (h *Handler) refreshSessionByRefreshTokenCookie(c *gin.Context) error {
	refreshToken, err := h.Cookie(c, refreshTokenCookieName)
	if err != nil {
		return err
	}

	newAccessToken, newRefreshToken, err := h.svc.UserAuthorization.RefreshSession(refreshToken)
	if err != nil {
		return err
	}

	accessTokenClaims, err := h.svc.UserAuthorization.GetAccessTokenClaims(newAccessToken)
	if err != nil {
		return err
	}

	h.setTokensCookies(c, newAccessToken, newRefreshToken)

	return h.validateAndSetUserIDForContext(c, accessTokenClaims.Subject)
}

func setHandlerNameToLogEntry(c *gin.Context, handlerName string) {
	logEntryValue, ok := c.Get(ctxLogEntry)
	if !ok {
		return
	}

	logEntry, ok := logEntryValue.(*logrus.Entry)
	if !ok {
		return
	}

	logEntry.WithField("method", handlerName)
	c.Set(ctxLogEntry, logEntry)
}

func (h *Handler) getLogEntry(c *gin.Context) *logrus.Entry {
	logEntryValue, ok := c.Get(ctxLogEntry)
	if !ok {
		return logrus.NewEntry(h.log)
	}

	logEntry, ok := logEntryValue.(*logrus.Entry)
	if !ok {
		return logrus.NewEntry(h.log)
	}

	return logEntry
}

func (h *Handler) validateAndSetUserIDForContext(c *gin.Context, userIDStr string) error {
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return err
	}

	if err = h.validateUserID(c, userID); err != nil {
		return err
	}

	c.Set(ctxUserID, userID)

	return nil
}

func (h *Handler) validateUserID(c *gin.Context, userID uint64) error {
	user, err := h.svc.User.GetUserByID(c, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	return nil
}

func getUserIDFromContext(c *gin.Context) (uint64, error) {
	userIDValue, ok := c.Get(ctxUserID)
	if !ok {
		return 0, errors.New("failed to get user id from context")
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		return 0, errors.Errorf("user id from context has not valid type: %T", userIDValue)
	}

	return userID, nil
}

// NOT USEFUL CODE BELOW.

// validateTokenHeader gets accessToken from header and validate it.
// on success it puts accessToken data to ctx and returns nil. else it returns error.
func (h *Handler) validateTokenHeader(c *gin.Context) error {
	header := c.GetHeader("Authorization")
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return ErrNotValidAuthorizationHeader
	}

	accessToken := headerParts[1]
	accessTokenClaims, err := h.svc.UserAuthorization.ValidateAccessToken(accessToken)
	if err != nil {
		return err
	}

	return h.validateAndSetUserIDForContext(c, accessTokenClaims.Subject)
}
