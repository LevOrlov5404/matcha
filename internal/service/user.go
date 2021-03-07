package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/LevOrlov5404/matcha/internal/repository"
	"github.com/LevOrlov5404/matcha/models"
)

const (
	tokenTTL = 12 * time.Hour
)

type (
	UserService struct {
		repo       repository.User
		salt       string
		signingKey string
	}

	tokenClaims struct {
		jwt.StandardClaims
		UserID int64 `json:"userId"`
	}
)

func NewUserService(repo repository.User, salt, signingKey string) *UserService {
	return &UserService{
		repo:       repo,
		salt:       salt,
		signingKey: signingKey,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user models.UserToCreate) (int64, error) {
	user.Password = generatePasswordHash(user.Password, s.salt)
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) GetUserByEmailPassword(ctx context.Context, email, password string) (models.UserToGet, error) {
	return s.repo.GetUserByEmailPassword(ctx, email, generatePasswordHash(password, s.salt))
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (models.UserToGet, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, user models.UserToCreate) error {
	user.Password = generatePasswordHash(user.Password, s.salt)
	return s.repo.UpdateUser(ctx, id, user)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.UserToGet, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) GenerateToken(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmailPassword(ctx, email, generatePasswordHash(password, s.salt))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.ID,
	})

	return token.SignedString([]byte(s.signingKey))
}

func (s *UserService) ParseToken(accessToken string) (int64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}

func generatePasswordHash(password, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
