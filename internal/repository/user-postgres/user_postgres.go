package user_postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	iErrs "github.com/LevOrlov5404/matcha/internal/errors"
	"github.com/LevOrlov5404/matcha/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	usersTable = "users"
)

type UserPostgres struct {
	db        *sqlx.DB
	dbTimeout time.Duration
}

func NewUserPostgres(db *sqlx.DB, dbTimeout time.Duration) *UserPostgres {
	return &UserPostgres{
		db:        db,
		dbTimeout: dbTimeout,
	}
}

func (r *UserPostgres) CreateUser(ctx context.Context, user models.UserToCreate) (uint64, error) {
	query := fmt.Sprintf(`
INSERT INTO %s (email, username, first_name, last_name, password)
values ($1, $2, $3, $4, $5) RETURNING id`, usersTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(dbCtx, query,
		user.Email, user.Username, user.FirstName, user.LastName, user.Password)
	if err := row.Err(); err != nil {
		return 0, getDBError(err)
	}

	var id uint64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserPostgres) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := fmt.Sprintf(`
SELECT id, email, username, first_name, last_name, password FROM %s WHERE username=$1`, usersTable)
	var user models.User

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if err := r.db.GetContext(dbCtx, &user, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
	}

	return &user, nil
}

func (r *UserPostgres) GetUserByID(ctx context.Context, id uint64) (*models.User, error) {
	query := fmt.Sprintf(`
SELECT id, email, username, first_name, last_name, password, is_email_confirmed FROM %s WHERE id=$1`, usersTable)
	var user models.User

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	if err := r.db.GetContext(dbCtx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
	}

	return &user, nil
}

func (r *UserPostgres) UpdateUser(ctx context.Context, user models.User) error {
	query := fmt.Sprintf(`
UPDATE %s SET email = $1, username = $2, first_name = $3, last_name = $4, password = $5 WHERE id = $6`, usersTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, user.Email, user.Username, user.FirstName, user.LastName, user.Password, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserPostgres) GetAllUsers(ctx context.Context) ([]models.User, error) {
	query := fmt.Sprintf(`
SELECT id, email, username, first_name, last_name, is_email_confirmed FROM %s`, usersTable)
	var users []models.User

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &users, query)

	return users, err
}

func (r *UserPostgres) DeleteUser(ctx context.Context, id uint64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, usersTable)

	dbCtx, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func getDBError(err error) error {
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Class() < "50" { // business error
			return iErrs.NewBusiness(err, err.Detail)
		}

		return iErrs.New(err)
	}

	return err
}
