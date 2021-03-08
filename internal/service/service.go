package service

import (
	"context"
	"time"

	"github.com/LevOrlov5404/matcha/internal/models"
	"github.com/LevOrlov5404/matcha/internal/repository"
)

type (
	User interface {
		CreateUser(ctx context.Context, user models.UserToCreate) (int64, error)
		GetUserByEmailPassword(ctx context.Context, email, password string) (models.UserToGet, error)
		GetUserByID(ctx context.Context, id int64) (models.UserToGet, error)
		UpdateUser(ctx context.Context, id int64, user models.UserToCreate) error
		GetAllUsers(ctx context.Context) ([]models.UserToGet, error)
		DeleteUser(ctx context.Context, id int64) error
		GenerateToken(ctx context.Context, email, password string) (string, error)
		ParseToken(token string) (int64, error)
	}

	Service struct {
		User
	}

	Options struct {
		TokenLifetime    time.Duration
		SigningKey       string
		UserPasswordSalt string
	}
)

func NewService(repo *repository.Repository, options Options) *Service {
	return &Service{
		User: NewUserService(
			repo.User, options.TokenLifetime, options.SigningKey, options.UserPasswordSalt,
		),
	}
}
