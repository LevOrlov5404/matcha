package service

import (
	"context"
	"errors"
	"time"

	iErrs "github.com/l-orlov/matcha/internal/errors"
	"github.com/l-orlov/matcha/internal/models"
	"github.com/l-orlov/matcha/internal/repository"
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

	hashedPassword, err := models.HashPassword(user.Password)
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
	return s.repo.UpdateUser(ctx, user)
}

func (s *UserService) SetUserPassword(ctx context.Context, userID uint64, password string) error {
	hashedPassword, err := models.HashPassword(password)
	if err != nil {
		return iErrs.New(err)
	}

	return s.repo.UpdateUserPassword(ctx, userID, hashedPassword)
}

func (s *UserService) ChangeUserPassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return iErrs.NewBusiness(errors.New("user does not exist"), "")
	}

	if !models.CheckPasswordHash(user.Password, oldPassword) {
		return iErrs.NewBusiness(errors.New("wrong password"), "")
	}

	hashedPassword, err := models.HashPassword(newPassword)
	if err != nil {
		return iErrs.New(err)
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

func (s *UserService) GetUserProfileByID(ctx context.Context, id uint64) (*models.UserProfile, error) {
	return s.repo.GetUserProfileByID(ctx, id)
}

func (s *UserService) UpdateUserProfile(ctx context.Context, user models.UserProfile) error {
	return s.repo.UpdateUserProfile(ctx, user)
}
