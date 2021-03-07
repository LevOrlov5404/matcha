package repository

import (
	"context"
	"fmt"

	"github.com/LevOrlov5404/matcha/models"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) CreateUser(ctx context.Context, user models.UserToCreate) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, surname, email, password_hash) values ($1, $2, $3, $4) RETURNING id", usersTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(dbCtx, query, user.Name, user.Surname, user.Email, user.Password)
	if err := row.Err(); err != nil {
		return 0, err
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserPostgres) GetUserByEmailPassword(ctx context.Context, email, password string) (models.UserToGet, error) {
	query := fmt.Sprintf("SELECT id, name, surname, email FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	var user models.UserToGet

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.GetContext(dbCtx, &user, query, email, password)

	return user, err
}

func (r *UserPostgres) GetUserByID(ctx context.Context, id int64) (models.UserToGet, error) {
	query := fmt.Sprintf("SELECT id, name, surname, email FROM %s WHERE id=$1", usersTable)
	var user models.UserToGet

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.GetContext(dbCtx, &user, query, id)

	return user, err
}

func (r *UserPostgres) UpdateUser(ctx context.Context, id int64, user models.UserToCreate) error {
	query := fmt.Sprintf("UPDATE %s SET name = $1, surname = $2, email = $3, password_hash = $4 WHERE id = $5", usersTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, user.Name, user.Surname, user.Email, user.Password, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserPostgres) GetAllUsers(ctx context.Context) ([]models.UserToGet, error) {
	query := fmt.Sprintf("SELECT id, name, surname, email FROM %s", usersTable)
	var users []models.UserToGet

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	err := r.db.SelectContext(dbCtx, &users, query)

	return users, err
}

func (r *UserPostgres) DeleteUser(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", usersTable)

	dbCtx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	_, err := r.db.ExecContext(dbCtx, query, id)
	if err != nil {
		return err
	}

	return nil
}
