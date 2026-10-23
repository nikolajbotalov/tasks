package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

func RunMigrations(ctx context.Context, db *sql.DB, migrationsFS fs.FS) error {
	source, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("migrations source: %w", err)
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("migration connection: %w", err)
	}

	driver, err := postgres.WithConnection(ctx, conn, &postgres.Config{})
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("create psql migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		_ = driver.Close()
		return fmt.Errorf("create migrate: %w", err)
	}

	upErr := m.Up()
	sourceErr, dbErr := m.Close()

	switch {
	case upErr != nil && !errors.Is(upErr, migrate.ErrNoChange):
		return fmt.Errorf("migrate up: %w", upErr)
	case dbErr != nil:
		return fmt.Errorf("close migration db: %w", dbErr)
	case sourceErr != nil:
		return fmt.Errorf("close migration source: %w", sourceErr)
	}

	return nil
}
