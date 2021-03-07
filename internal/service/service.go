package service

import (
	"context"

	"github.com/LevOrlov5404/matcha/internal/repository"
	"github.com/LevOrlov5404/matcha/models"
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
)

func NewService(repo *repository.Repository, salt, signingKey string) *Service {
	return &Service{
		User: NewUserService(repo.User, salt, signingKey),
	}
}
