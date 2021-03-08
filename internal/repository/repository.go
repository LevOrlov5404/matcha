package repository

import (
	"context"

	"github.com/LevOrlov5404/matcha/internal/models"
	"github.com/jmoiron/sqlx"
)

type (
	User interface {
		CreateUser(ctx context.Context, user models.UserToCreate) (int64, error)
		GetUserByEmailPassword(ctx context.Context, email, password string) (models.UserToGet, error)
		GetUserByID(ctx context.Context, id int64) (models.UserToGet, error)
		UpdateUser(ctx context.Context, id int64, user models.UserToCreate) error
		GetAllUsers(ctx context.Context) ([]models.UserToGet, error)
		DeleteUser(ctx context.Context, id int64) error
	}

	Repository struct {
		User
	}
)

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserPostgres(db),
	}
}
