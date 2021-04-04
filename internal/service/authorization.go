package service

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/l-orlov/matcha/internal/config"
	"github.com/l-orlov/matcha/internal/models"
	"github.com/l-orlov/matcha/internal/repository"
	"github.com/pkg/errors"
)

var (
	ErrNotActiveAccessToken = errors.New("not active accessToken")
	ErrSessionNotFound      = errors.New("session not found")
)

type (
	AuthorizationService struct {
		cfg  *config.Config
		repo repository.SessionCache
	}
)

func NewAuthorizationService(cfg *config.Config, repo *repository.Repository) *AuthorizationService {
	return &AuthorizationService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *AuthorizationService) CreateSession(userID string) (accessToken, refreshToken string, err error) {
	accessTokenID := uuid.New().String()
	accessToken, err = newToken(
		userID, accessTokenID, s.cfg.JWT.SigningKey, s.cfg.JWT.AccessTokenLifetime.Duration(),
	)
	if err != nil {
		return "", "", err
	}

	refreshToken = uuid.New().String()

	err = s.repo.PutSessionAndAccessToken(models.Session{
		UserID:        userID,
		AccessTokenID: accessTokenID,
	}, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthorizationService) ValidateAccessToken(accessToken string) (*jwt.StandardClaims, error) {
	accessTokenClaims, err := validateToken(accessToken, s.cfg.JWT.SigningKey)
	if err != nil {
		return nil, err
	}

	// check accessToken is active
	if _, err := s.repo.GetAccessTokenData(accessTokenClaims.Id); err != nil {
		if errors.Is(err, redis.ErrNil) {
			return nil, ErrNotActiveAccessToken
		}

		return nil, err
	}

	return accessTokenClaims, nil
}

func (s *AuthorizationService) RefreshSession(
	currentRefreshToken string,
) (accessToken, refreshToken string, err error) {
	session, err := s.repo.GetSession(currentRefreshToken)
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return "", "", ErrSessionNotFound
		}

		return "", "", err
	}

	if err = s.repo.DeleteSession(currentRefreshToken); err != nil {
		return "", "", err
	}

	if err = s.repo.DeleteUserToSession(session.UserID, currentRefreshToken); err != nil {
		return "", "", err
	}

	if err = s.repo.DeleteAccessToken(session.AccessTokenID); err != nil {
		return "", "", err
	}

	return s.CreateSession(session.UserID)
}

func (s *AuthorizationService) RevokeSession(accessToken string) error {
	accessTokenClaims, err := validateToken(accessToken, s.cfg.JWT.SigningKey)
	if err != nil {
		return err
	}

	refreshToken, err := s.repo.GetAccessTokenData(accessTokenClaims.Id)
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return ErrNotActiveAccessToken
		}

		return err
	}

	if err := s.repo.DeleteAccessToken(accessTokenClaims.Id); err != nil {
		return err
	}

	session, _ := s.repo.GetSession(refreshToken)
	if session != nil {
		if err = s.repo.DeleteUserToSession(session.UserID, refreshToken); err != nil {
			return err
		}
	}

	if err := s.repo.DeleteSession(refreshToken); err != nil {
		return err
	}

	return nil
}

func (s *AuthorizationService) GetAccessTokenClaims(accessToken string) (*jwt.StandardClaims, error) {
	return getTokenClaims(accessToken, s.cfg.JWT.SigningKey)
}

func newToken(userID, tokenID, signingKey string, lifetime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Id:        tokenID,
		NotBefore: time.Now().Unix(),
		ExpiresAt: time.Now().Add(lifetime).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   userID,
	})

	return token.SignedString([]byte(signingKey))
}

func validateToken(token string, signingKey string) (*jwt.StandardClaims, error) {
	claims, err := getTokenClaims(token, signingKey)
	if err != nil {
		return nil, errors.Wrap(err, "not valid token")
	}

	// check accessToken has not expired
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("expired accessToken")
	}

	return claims, nil
}

func getTokenClaims(tokenString string, signingKey string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	return token.Claims.(*jwt.StandardClaims), nil
}
