package repository

import (
	"fmt"

	customErrs "github.com/LevOrlov5404/matcha/internal/custom-errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	usersTable = "users"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func ConnectToDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", initConnectionString(cfg))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func initConnectionString(cfg Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)
}

func getDBError(err error) error {
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Class() < "50" { // business error
			return customErrs.NewBusiness(err, err.Detail)
		}

		return customErrs.New(err)
	}

	return err
}
