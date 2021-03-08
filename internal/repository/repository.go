package repository

import (
	"context"
	"time"

	"github.com/LevOrlov5404/matcha/internal/models"
	"github.com/jmoiron/sqlx"
)

type (
	User interface {
		CreateUser(ctx context.Context, user models.UserToCreate) (uint64, error)
		GetUserByUsername(ctx context.Context, username string) (*models.User, error)
		GetUserByID(ctx context.Context, id uint64) (*models.User, error)
		UpdateUser(ctx context.Context, user models.User) error
		GetAllUsers(ctx context.Context) ([]models.User, error)
		DeleteUser(ctx context.Context, id uint64) error
	}

	Repository struct {
		User
	}
)

func NewRepository(db *sqlx.DB, dbTimeout time.Duration) *Repository {
	return &Repository{
		User: NewUserPostgres(db, dbTimeout),
	}
}
