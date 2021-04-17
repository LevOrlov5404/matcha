package service

import (
	"context"
	"time"

	ierrors "github.com/l-orlov/matcha/internal/errors"
	"github.com/l-orlov/matcha/internal/models"
	"github.com/l-orlov/matcha/internal/repository"
	"github.com/pkg/errors"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUsernameIsTaken = errors.New("username is already taken")
	ErrEmailIsTaken    = errors.New("user with this email already exists")
	ErrWrongPassword   = errors.New("wrong password")
)

type (
	UserService struct {
		repo                repository.User
		accessTokenLifetime time.Duration
	}
)

func NewUserService(
	repo repository.User, tokenLifetime time.Duration,
) *UserService {
	return &UserService{
		repo:                repo,
		accessTokenLifetime: tokenLifetime,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user models.UserToCreate) (uint64, error) {
	existingUser, err := s.repo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return 0, err
	}

	if existingUser != nil {
		return 0, ierrors.NewBusiness(ErrUsernameIsTaken, "")
	}

	existingUser, err = s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return 0, err
	}

	if existingUser != nil {
		return 0, ierrors.NewBusiness(ErrEmailIsTaken, "")
	}

	hashedPassword, err := models.HashPassword(user.Password)
	if err != nil {
		return 0, ierrors.New(err)
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
	return s.repo.UpdateUser(ctx, user)
}

func (s *UserService) SetUserPassword(ctx context.Context, userID uint64, password string) error {
	hashedPassword, err := models.HashPassword(password)
	if err != nil {
		return ierrors.New(err)
	}

	return s.repo.UpdateUserPassword(ctx, userID, hashedPassword)
}

func (s *UserService) ChangeUserPassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return ierrors.NewBusiness(ErrUserNotFound, "")
	}

	if !models.CheckPasswordHash(user.Password, oldPassword) {
		return ierrors.NewBusiness(ErrWrongPassword, "")
	}

	hashedPassword, err := models.HashPassword(newPassword)
	if err != nil {
		return ierrors.New(err)
	}

	return s.repo.UpdateUserPassword(ctx, userID, hashedPassword)
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
