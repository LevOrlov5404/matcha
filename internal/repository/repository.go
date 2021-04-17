package repository

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/l-orlov/matcha/internal/config"
	"github.com/l-orlov/matcha/internal/models"
	cacheredis "github.com/l-orlov/matcha/internal/repository/cache-redis"
	storageminio "github.com/l-orlov/matcha/internal/repository/storage-minio"
	userpostgres "github.com/l-orlov/matcha/internal/repository/user-postgres"
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
		GetUserProfileByID(ctx context.Context, id uint64) (*models.UserProfile, error)
		UpdateUserProfile(ctx context.Context, user models.UserProfile) error
		UpdateUserAvatarPath(ctx context.Context, userID uint64, avatarPath string) error
	}
	UserPictures interface {
		CreateUserPicture(ctx context.Context, picture models.UserPicture) error
		GetUserPictureByUUID(ctx context.Context, uuid uuid.UUID) (*models.UserPicture, error)
		GetUserPicturesByUserID(ctx context.Context, userID uint64) ([]models.UserPicture, error)
		DeleteUserPicture(ctx context.Context, uuid uuid.UUID) error
	}
	SessionCache interface {
		PutSessionAndAccessToken(session models.Session, refreshToken string) error
		GetSession(refreshToken string) (*models.Session, error)
		DeleteSession(refreshToken string) error
		DeleteUserToSession(userID, refreshToken string) error
		GetAccessTokenData(accessTokenID string) (refreshToken string, err error)
		DeleteAccessToken(accessTokenID string) error
		AddUserBlocking(fingerprint string) (int64, error)
		GetUserBlocking(fingerprint string) (int, error)
		DeleteUserBlocking(fingerprint string) error
	}
	VerificationCache interface {
		PutEmailConfirmToken(userID uint64, token string) error
		GetEmailConfirmTokenData(token string) (userID uint64, err error)
		DeleteEmailConfirmToken(token string) error
		PutPasswordResetConfirmToken(userID uint64, token string) error
		GetPasswordResetConfirmTokenData(token string) (userID uint64, err error)
		DeletePasswordResetConfirmToken(token string) error
	}
	Storage interface {
		PutFile(ctx context.Context, bucketName, objectName, contentType string, reader io.Reader) error
		GetFileURL(ctx context.Context, bucket, objectName string, expires time.Duration) (url string, err error)
		DeleteFile(ctx context.Context, bucket, objectName string) error
	}
	Repository struct {
		User
		UserPictures
		SessionCache
		VerificationCache
		Storage
	}
)

func NewRepository(
	cfg *config.Config, log *logrus.Logger, db *sqlx.DB,
) (*Repository, error) {
	userRepo := userpostgres.NewUserPostgres(db, cfg.PostgresDB.Timeout.Duration())
	userPicturesRepo := userpostgres.NewUserPicturesPostgres(db, cfg.PostgresDB.Timeout.Duration())

	cacheLogEntry := logrus.NewEntry(log).WithFields(logrus.Fields{"source": "cache-redis"})
	cacheOptions := cacheredis.Options{
		AccessTokenLifetime:               int(cfg.JWT.AccessTokenLifetime.Duration().Seconds()),
		RefreshTokenLifetime:              int(cfg.JWT.RefreshTokenLifetime.Duration().Seconds()),
		UserBlockingLifetime:              int(cfg.UserBlocking.Lifetime.Duration().Seconds()),
		EmailConfirmTokenLifetime:         int(cfg.Verification.EmailConfirmTokenLifetime.Duration().Seconds()),
		PasswordResetConfirmTokenLifetime: int(cfg.Verification.PasswordResetConfirmTokenLifetime.Duration().Seconds()),
	}
	cache := cacheredis.New(cfg.Redis, cacheLogEntry, cacheOptions)

	storageEntry := logrus.NewEntry(log).WithFields(logrus.Fields{"source": "minio-storage"})
	storage, err := storageminio.New(storageminio.Config{
		Endpoint:  cfg.Minio.Endpoint.String(),
		AccessKey: cfg.Minio.AccessKey,
		SecretKey: cfg.Minio.SecretKey,
		UseSSL:    cfg.Minio.UseSSL,
		Timeout:   cfg.Minio.Timeout.Duration(),
	}, storageEntry)
	if err != nil {
		return nil, err
	}

	return &Repository{
		User:              userRepo,
		UserPictures:      userPicturesRepo,
		SessionCache:      cache,
		VerificationCache: cache,
		Storage:           storage,
	}, nil
}
