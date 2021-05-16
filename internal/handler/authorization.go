package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	ierrors "github.com/l-orlov/matcha/internal/errors"
	"github.com/l-orlov/matcha/internal/models"
	"github.com/pkg/errors"
)

const (
	accessTokenCookieName  = "_m_access"
	refreshTokenCookieName = "_m_refresh"
)

func (h *Handler) SignIn(c *gin.Context) {
	setHandlerNameToLogEntry(c, "SignIn")

	var user models.UserToSignIn
	if err := c.BindJSON(&user); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	userID, err := h.svc.AuthenticateUserByUsername(c, user.Username, user.Password, user.Fingerprint)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	accessToken, refreshToken, err := h.svc.CreateSession(strconv.FormatUint(userID, 10))
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	h.setTokensCookies(c, accessToken, refreshToken)
	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) ValidateAccessToken(c *gin.Context) {
	setHandlerNameToLogEntry(c, "ValidateAccessToken")

	var req models.ValidateAccessTokenRequest
	if err := c.BindJSON(&req); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if _, err := h.svc.UserAuthorization.ValidateAccessToken(req.AccessToken); err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) RefreshSession(c *gin.Context) {
	setHandlerNameToLogEntry(c, "RefreshSession")

	var req models.RefreshSessionRequest
	if err := c.BindJSON(&req); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	accessToken, refreshToken, err := h.svc.RefreshSession(req.RefreshToken)
	if err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	setHandlerNameToLogEntry(c, "Logout")

	accessToken, err := h.Cookie(c, accessTokenCookieName)
	if err != nil {
		h.getLogEntry(c).Debug(err)
	}

	if accessToken == "" {
		// try to get accessToken from request
		var req models.LogoutRequest
		if err := c.BindJSON(&req); err != nil {
			h.newErrorResponse(c, http.StatusBadRequest, err)
			return
		}

		accessToken = req.AccessToken
	}

	if err = h.svc.RevokeSession(accessToken); err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) ResetPassword(c *gin.Context) {
	setHandlerNameToLogEntry(c, "ResetPassword")

	email, ok := c.GetQuery("email")
	if !ok || email == "" {
		h.newErrorResponse(
			c, http.StatusBadRequest, ierrors.NewBusiness(ErrEmptyEmailParameter, ""),
		)
		return
	}

	user, err := h.svc.User.GetUserByEmail(c, email)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	if user == nil {
		h.newErrorResponse(
			c, http.StatusBadRequest, ierrors.NewBusiness(ErrUserNotFound, ""),
		)
		return
	}

	passwordResetConfirmToken, err := h.svc.Verification.CreatePasswordResetConfirmToken(user.ID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	// send token by email
	h.svc.Mailer.SendResetPasswordConfirm(user.Email, passwordResetConfirmToken)

	c.Status(http.StatusOK)
}

func (h *Handler) setTokensCookies(c *gin.Context, accessToken, refreshToken string) {
	if encodedAccessToken, err := h.options.SecureCookie.Encode(accessTokenCookieName, accessToken); err == nil {
		c.SetCookie(
			accessTokenCookieName, encodedAccessToken, h.options.AccessTokenCookieMaxAge,
			"/", h.cfg.Cookie.Domain, false, true,
		)
	} else {
		h.getLogEntry(c).Error(err)
	}

	if encodedRefreshToken, err := h.options.SecureCookie.Encode(refreshTokenCookieName, refreshToken); err == nil {
		c.SetCookie(
			refreshTokenCookieName, encodedRefreshToken, h.options.RefreshTokenCookieMaxAge,
			"/", h.cfg.Cookie.Domain, false, true,
		)
	} else {
		h.getLogEntry(c).Error(err)
	}
}

func (h *Handler) Cookie(c *gin.Context, name string) (string, error) {
	value, err := c.Cookie(name)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get cookie %s", name)
	}

	var decodedValue string
	if err := h.options.SecureCookie.Decode(name, value, &decodedValue); err != nil {
		return "", err
	}

	return decodedValue, nil
}
