package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"

	dbTimeout = 3 * time.Second
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
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
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
}
