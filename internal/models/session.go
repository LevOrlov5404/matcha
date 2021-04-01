package models

type Session struct {
	UserID        string `json:"userId"`
	AccessTokenID string `json:"accessTokenId"`
	Fingerprint   string `json:"fingerprint"`
}

type ValidateAccessTokenRequest struct {
	AccessToken string `json:"accessToken" binding:"required"`
}

type RefreshSessionRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
	Fingerprint  string `json:"fingerprint" binding:"required"`
}

type LogoutRequest struct {
	AccessToken string `json:"accessToken" binding:"required"`
}
