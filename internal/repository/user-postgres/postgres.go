package user_postgres

import (
	"fmt"

	"github.com/LevOrlov5404/matcha/internal/config"
	"github.com/jmoiron/sqlx"
)

func ConnectToDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", initConnectionString(cfg.PostgresDB))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func initConnectionString(cfg config.PostgresDB) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Address.Host, cfg.Address.Port, cfg.User, cfg.Password, cfg.Database)
}
