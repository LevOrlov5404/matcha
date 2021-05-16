package postgres

import (
	"database/sql"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/l-orlov/matcha/internal/config"
	"github.com/pkg/errors"
)

func MigrateSchema(db *sql.DB, cfg config.PostgresDB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+cfg.MigrationDir, cfg.Database, driver)
	if err != nil {
		return err
	}

	if err = m.Steps(1); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	return nil
}
