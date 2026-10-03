package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

func TableExists(ctx context.Context, db *sql.DB, tableName string) (bool, error) {
	var exists bool

	err := db.QueryRowContext(ctx, "SELECT to_regclass($1) IS NOT NULL", tableName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func EnsureSchema(ctx context.Context, db *sql.DB, tables ...string) error {
	for _, t := range tables {
		ok, err := TableExists(ctx, db, t)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("table %s does not exist", t)
		}
	}

	return nil
}
