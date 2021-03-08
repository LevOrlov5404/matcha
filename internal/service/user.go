package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/LevOrlov5404/matcha/internal/models"
	"github.com/LevOrlov5404/matcha/internal/repository"
	"github.com/dgrijalva/jwt-go"
)

type (
	UserService struct {
		repo          repository.User
		tokenLifetime time.Duration
		signingKey    string
		passwordSalt  string
	}
)

func NewUserService(
	repo repository.User, tokenLifetime time.Duration, signingKey, passwordSalt string,
) *UserService {
	return &UserService{
		repo:          repo,
		tokenLifetime: tokenLifetime,
		signingKey:    signingKey,
		passwordSalt:  passwordSalt,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user models.UserToCreate) (int64, error) {
	user.Password = generatePasswordHash(user.Password, s.passwordSalt)
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) GetUserByEmailPassword(ctx context.Context, email, password string) (models.UserToGet, error) {
	return s.repo.GetUserByEmailPassword(ctx, email, generatePasswordHash(password, s.passwordSalt))
}

func (s *UserService) GetUserByID(ctx context.Context, id int64) (models.UserToGet, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, user models.UserToCreate) error {
	user.Password = generatePasswordHash(user.Password, s.passwordSalt)
	return s.repo.UpdateUser(ctx, id, user)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.UserToGet, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) GenerateToken(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmailPassword(ctx, email, generatePasswordHash(password, s.passwordSalt))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(s.tokenLifetime).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.FormatInt(user.ID, 10),
	})

	return token.SignedString([]byte(s.signingKey))
}

func (s *UserService) ParseToken(accessToken string) (int64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims := token.Claims.(*jwt.StandardClaims)

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to get userID from token: %v", err)
	}

	return userID, nil
}

func generatePasswordHash(password, salt string) string {
	hash := sha1.New() // ToDo: change to secure hash
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
