package service

import (
	"context"

	"github.com/LevOrlov5404/matcha/internal/config"
	"github.com/LevOrlov5404/matcha/internal/models"
	"github.com/LevOrlov5404/matcha/internal/repository"
	"github.com/sirupsen/logrus"
)

type (
	RandomTokenGenerator interface {
		Generate(length, digitsNum, symbolsNum int, noUpper, allowRepeat bool) (string, error)
	}
	User interface {
		CreateUser(ctx context.Context, user models.UserToCreate) (uint64, error)
		GetUserByID(ctx context.Context, id uint64) (*models.User, error)
		GetUserByEmail(ctx context.Context, email string) (*models.User, error)
		UpdateUser(ctx context.Context, user models.User) error
		SetUserPassword(ctx context.Context, userID uint64, password string) error
		ChangeUserPassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error
		GetAllUsers(ctx context.Context) ([]models.User, error)
		DeleteUser(ctx context.Context, id uint64) error
		ConfirmEmail(ctx context.Context, id uint64) error
		GenerateToken(ctx context.Context, username, password string) (string, error)
		ParseToken(token string) (uint64, error)
	}
	Verification interface {
		CreateEmailConfirmToken(userID uint64) (string, error)
		VerifyEmailConfirmToken(emailConfirmToken string) (userID uint64, err error)
		CreatePasswordResetConfirmToken(userID uint64) (string, error)
		VerifyPasswordResetConfirmToken(confirmToken string) (userID uint64, err error)
	}
	Mailer interface {
		SendEmailConfirm(toEmail, token string)
		SendResetPasswordConfirm(toEmail, token string)
	}
	Service struct {
		User
		Verification
		Mailer
	}
)

func NewService(
	cfg *config.Config, log *logrus.Logger,
	repo *repository.Repository, generator RandomTokenGenerator,
	mailer Mailer,
) *Service {
	verificationLogEntry := logrus.NewEntry(log).WithFields(logrus.Fields{"source": "verificationService"})

	return &Service{
		User: NewUserService(
			repo.User, cfg.JWT.AccessTokenLifetime.Duration(), cfg.JWT.SigningKey,
		),
		Verification: NewVerificationService(verificationLogEntry, repo.Cache, generator),
		Mailer:       mailer,
	}
}
