package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	iErrs "github.com/LevOrlov5404/matcha/internal/errors"
	"github.com/LevOrlov5404/matcha/internal/models"
	"github.com/LevOrlov5404/matcha/internal/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserService struct {
		repo                repository.User
		accessTokenLifetime time.Duration
		signingKey          string
	}
)

func NewUserService(
	repo repository.User, tokenLifetime time.Duration, signingKey string,
) *UserService {
	return &UserService{
		repo:                repo,
		accessTokenLifetime: tokenLifetime,
		signingKey:          signingKey,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user models.UserToCreate) (uint64, error) {
	existingUser, err := s.repo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return 0, err
	}

	if existingUser != nil {
		return 0, iErrs.NewBusiness(errors.New("username is already taken"), "")
	}

	existingUser, err = s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return 0, err
	}

	if existingUser != nil {
		return 0, iErrs.NewBusiness(errors.New("user with this email already exists"), "")
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return 0, iErrs.New(err)
	}

	user.Password = hashedPassword

	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id uint64) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *UserService) UpdateUser(ctx context.Context, user models.User) error {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return iErrs.New(err)
	}

	user.Password = hashedPassword

	return s.repo.UpdateUser(ctx, user)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAllUsers(ctx)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint64) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *UserService) ConfirmEmail(ctx context.Context, id uint64) error {
	return s.repo.ConfirmEmail(ctx, id)
}

func (s *UserService) GenerateToken(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", iErrs.NewBusiness(errors.New("user not found"), "")
	}

	if !CheckPasswordHash(user.Password, password) {
		return "", iErrs.NewBusiness(errors.New("wrong password"), "")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(s.accessTokenLifetime).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   strconv.FormatUint(user.ID, 10),
	})

	return token.SignedString([]byte(s.signingKey))
}

func (s *UserService) ParseToken(accessToken string) (uint64, error) {
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

	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to get userID from token: %v", err)
	}

	return userID, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
