package service

import (
	"context"
	"io"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
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
	UserProfile interface {
		GetUserProfileByID(ctx context.Context, id uint64) (*models.UserProfile, error)
		UpdateUserProfile(ctx context.Context, user models.UserProfile) error
		UploadUserAvatar(ctx context.Context, userID uint64, file io.ReadSeeker) error
		DeleteUserAvatar(ctx context.Context, userID uint64) error
		UploadUserPicture(ctx context.Context, userID uint64, file io.ReadSeeker) error
		GetUserPicturesByUserID(ctx context.Context, userID uint64) ([]models.UserPicture, error)
		DeleteUserPicture(ctx context.Context, uuid uuid.UUID) error
	}
	Service struct {
		User
		UserAuthentication
		UserAuthorization
		Verification
		Mailer
		UserProfile
	}
)

func NewService(
	cfg *config.Config, log *logrus.Logger,
	repo *repository.Repository, generator RandomTokenGenerator,
	mailer Mailer,
) *Service {
	authenticationLogEntry := logrus.NewEntry(log).WithFields(logrus.Fields{"source": "authentication-svc"})
	verificationLogEntry := logrus.NewEntry(log).WithFields(logrus.Fields{"source": "verification-svc"})
	profileLogEntry := logrus.NewEntry(log).WithFields(logrus.Fields{"source": "user-profile-svc"})

	return &Service{
		User:               NewUserService(repo.User, cfg.JWT.AccessTokenLifetime.Duration()),
		UserAuthentication: NewAuthenticationService(cfg, authenticationLogEntry, repo),
		UserAuthorization:  NewAuthorizationService(cfg, repo),
		Verification:       NewVerificationService(verificationLogEntry, repo.VerificationCache, generator),
		Mailer:             mailer,
		UserProfile:        NewUserProfileService(profileLogEntry, cfg.MaxUserPicturesNum, cfg.FilePathTemplates, repo),
	}
}
