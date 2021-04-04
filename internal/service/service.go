package service

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/l-orlov/matcha/internal/config"
	"github.com/l-orlov/matcha/internal/models"
	"github.com/l-orlov/matcha/internal/repository"
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
		GetUserProfileByID(ctx context.Context, id uint64) (*models.UserProfile, error)
		UpdateUserProfile(ctx context.Context, user models.UserProfile) error
	}
	UserAuthentication interface {
		AuthenticateUserByUsername(ctx context.Context, username, password, fingerprint string) (userID uint64, err error)
	}
	UserAuthorization interface {
		CreateSession(userID string) (accessToken, refreshToken string, err error)
		ValidateAccessToken(accessToken string) (*jwt.StandardClaims, error)
		RefreshSession(currentRefreshToken string) (accessToken, refreshToken string, err error)
		RevokeSession(accessToken string) error
		GetAccessTokenClaims(accessToken string) (*jwt.StandardClaims, error)
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
		UserAuthentication
		UserAuthorization
		Verification
		Mailer
	}
)

func NewService(
	cfg *config.Config, log *logrus.Logger,
	repo *repository.Repository, generator RandomTokenGenerator,
	mailer Mailer,
) *Service {
	authenticationLogEntry := logrus.NewEntry(log).WithFields(logrus.Fields{"source": "authenticationService"})
	verificationLogEntry := logrus.NewEntry(log).WithFields(logrus.Fields{"source": "verificationService"})

	return &Service{
		User:               NewUserService(repo.User, cfg.JWT.AccessTokenLifetime.Duration(), cfg.JWT.SigningKey),
		UserAuthentication: NewAuthenticationService(cfg, authenticationLogEntry, repo),
		UserAuthorization:  NewAuthorizationService(cfg, repo),
		Verification:       NewVerificationService(verificationLogEntry, repo.VerificationCache, generator),
		Mailer:             mailer,
	}
}
