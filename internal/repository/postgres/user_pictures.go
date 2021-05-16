package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/l-orlov/matcha/internal/models"
	"github.com/pkg/errors"
)

const (
	usersPicturesTable = "users_pictures"
)

type UserPicturesPostgres struct {
	db        *sqlx.DB
	dbTimeout time.Duration
}

func NewUserPicturesPostgres(db *sqlx.DB, dbTimeout time.Duration) *UserPicturesPostgres {
	return &UserPicturesPostgres{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *UserPicturesPostgres) CreateUserPicture(ctx context.Context, picture models.UserPicture) error {
	query := fmt.Sprintf(`
INSERT INTO %s (uuid, user_id, picture_path) VALUES ($1, $2, $3)`, usersPicturesTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, &picture.UUID, &picture.UserID, &picture.PicturePath)
	if err != nil {
		return getDBError(err)
	}

	return nil
}

func (r *UserPicturesPostgres) GetUserPictureByUUID(ctx context.Context, uuid uuid.UUID) (*models.UserPicture, error) {
	query := fmt.Sprintf(`SELECT uuid, user_id, picture_path FROM %s WHERE uuid=$1`, usersPicturesTable)
	var picture models.UserPicture

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if err := r.db.GetContext(dbCtx, &picture, query, &uuid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &picture, nil
}

func (r *UserPicturesPostgres) GetUserPicturesByUserID(ctx context.Context, userID uint64) ([]models.UserPicture, error) {
	query := fmt.Sprintf(`SELECT uuid, user_id, picture_path FROM %s WHERE user_id=$1`, usersPicturesTable)
	var pictures []models.UserPicture

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &pictures, query, &userID)

	return pictures, err
}

func (r *UserPicturesPostgres) DeleteUserPicture(ctx context.Context, uuid uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE uuid = $1`, usersPicturesTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if _, err := r.db.ExecContext(dbCtx, query, &uuid); err != nil {
		return err
	}

	return nil
}
