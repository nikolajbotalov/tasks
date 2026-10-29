package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *authRepository) CreateUser(ctx context.Context, id, email, password string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin create user: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO users (id, email, password)
				VALUES ($1, $2, $3)
				ON CONFLICT (email)
				DO NOTHING RETURNING id`, id, email, password,
	).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrEmailExists
		}

		return fmt.Errorf("insert user: %w", err)
	}

	if _, err = tx.ExecContext(ctx,
		`INSERT INTO profiles (user_id)
					VALUES ($1)`, id); err != nil {
		return fmt.Errorf("insert profile for user: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit create user: %w", err)
	}

	return nil
}
