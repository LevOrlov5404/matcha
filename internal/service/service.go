package service

import (
	"context"
	"time"

	"github.com/LevOrlov5404/matcha/internal/models"
	"github.com/LevOrlov5404/matcha/internal/repository"
)

type (
	User interface {
		CreateUser(ctx context.Context, user models.UserToCreate) (uint64, error)
		GetUserByID(ctx context.Context, id uint64) (*models.User, error)
		UpdateUser(ctx context.Context, user models.User) error
		GetAllUsers(ctx context.Context) ([]models.User, error)
		DeleteUser(ctx context.Context, id uint64) error
		GenerateToken(ctx context.Context, username, password string) (string, error)
		ParseToken(token string) (uint64, error)
	}

	Service struct {
		User
	}

	Options struct {
		TokenLifetime time.Duration
		SigningKey    string
	}
)

func NewService(repo *repository.Repository, options Options) *Service {
	return &Service{
		User: NewUserService(repo.User, options.TokenLifetime, options.SigningKey),
	}
}
