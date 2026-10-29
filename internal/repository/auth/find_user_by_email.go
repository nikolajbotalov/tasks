package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *authRepository) FindUserByEmail(ctx context.Context, email string) (domain.UserRow, error) {
	var userRow domain.UserRow

	err := r.db.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = $1", email).
		Scan(&userRow.ID, &userRow.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userRow, domain.ErrUserNotFound
		}
		return userRow, fmt.Errorf("find user by email: %w", err)
	}

	return userRow, nil
}
