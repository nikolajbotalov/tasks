package postgres

import (
	"TaskFlow/internal/config"
	"TaskFlow/migrations"
	"context"
	"database/sql"
	"fmt"
	"time"
)

func ConnectDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.Database.DBHost, cfg.Database.DBPort, cfg.Database.DBName, cfg.Database.DBUser, cfg.Database.DBPassword, cfg.Database.DBSSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	db.SetMaxOpenConns(cfg.Database.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.DBMaxIdleConns)
	db.SetConnMaxLifetime(time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	if err = RunMigrations(context.Background(), db, migrations.FS); err != nil {
		return nil, fmt.Errorf("migrate db: %w", err)
	}

	return db, nil
}
