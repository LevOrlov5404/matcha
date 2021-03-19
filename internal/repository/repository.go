package repository

import (
	"context"

	"github.com/LevOrlov5404/matcha/internal/config"
	"github.com/LevOrlov5404/matcha/internal/models"
	cacheRedis "github.com/LevOrlov5404/matcha/internal/repository/cache-redis"
	userPostgres "github.com/LevOrlov5404/matcha/internal/repository/user-postgres"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type (
	User interface {
		CreateUser(ctx context.Context, user models.UserToCreate) (uint64, error)
		GetUserByUsername(ctx context.Context, username string) (*models.User, error)
		GetUserByID(ctx context.Context, id uint64) (*models.User, error)
		GetUserByEmail(ctx context.Context, email string) (*models.User, error)
		UpdateUser(ctx context.Context, user models.User) error
		UpdateUserPassword(ctx context.Context, userID uint64, password string) error
		GetAllUsers(ctx context.Context) ([]models.User, error)
		DeleteUser(ctx context.Context, id uint64) error
		ConfirmEmail(ctx context.Context, id uint64) error
	}
	Cache interface {
		PutEmailConfirmToken(userID uint64, token string) error
		GetEmailConfirmTokenData(token string) (userID uint64, err error)
		DeleteEmailConfirmToken(token string) error
		PutResetPasswordConfirmToken(userID uint64, token string) error
		GetResetPasswordConfirmTokenData(token string) (userID uint64, err error)
		DeleteResetPasswordConfirmToken(token string) error
	}
	Repository struct {
		User
		Cache
	}
)

func NewRepository(
	cfg *config.Config, log *logrus.Logger, db *sqlx.DB,
) *Repository {
	cacheEntry := logrus.NewEntry(log).WithFields(logrus.Fields{"source": "cacheRedis"})
	cacheOptions := cacheRedis.Options{
		EmailConfirmTokenLifetime:         int(cfg.Verification.EmailConfirmTokenLifetime.Duration().Seconds()),
		ResetPasswordConfirmTokenLifetime: int(cfg.Verification.ResetPasswordConfirmTokenLifetime.Duration().Seconds()),
	}

	return &Repository{
		User:  userPostgres.NewUserPostgres(db, cfg.PostgresDB.Timeout.Duration()),
		Cache: cacheRedis.New(cfg.Redis, cacheEntry, cacheOptions),
	}
}
