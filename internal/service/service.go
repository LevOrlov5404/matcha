package service

import (
	"context"

	"github.com/LevOrlov5404/matcha/internal/config"
	"github.com/LevOrlov5404/matcha/internal/models"
	"github.com/LevOrlov5404/matcha/internal/repository"
)

type (
	RandomTokenGenerator interface {
		Generate(length, numDigits, numSymbols int, noUpper, allowRepeat bool) (string, error)
	}
	User interface {
		CreateUser(ctx context.Context, user models.UserToCreate) (uint64, error)
		GetUserByID(ctx context.Context, id uint64) (*models.User, error)
		UpdateUser(ctx context.Context, user models.User) error
		GetAllUsers(ctx context.Context) ([]models.User, error)
		DeleteUser(ctx context.Context, id uint64) error
		ConfirmEmail(ctx context.Context, id uint64) error
		GenerateToken(ctx context.Context, username, password string) (string, error)
		ParseToken(token string) (uint64, error)
	}
	Verification interface {
		CreateEmailConfirmToken(userID uint64) (string, error)
		VerifyEmailConfirmToken(emailConfirmToken string) (userID uint64, err error)
	}
	Mailer interface {
		SendEmailConfirm(toEmail, token string) error
	}
	Service struct {
		User
		Verification
		Mailer
	}
)

func NewService(
	cfg *config.Config, repo *repository.Repository, generator RandomTokenGenerator,
) *Service {
	return &Service{
		User: NewUserService(
			repo.User, cfg.JWT.AccessTokenLifetime.Duration(), cfg.JWT.SigningKey,
		),
		Verification: NewVerificationService(repo.Cache, generator),
		Mailer:       NewMailerService(cfg.Mailer),
	}
}
