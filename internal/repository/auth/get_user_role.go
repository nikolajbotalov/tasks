package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *authRepository) GetUserRole(ctx context.Context, id uuid.UUID) (domain.Role, error) {
	var role string

	err := r.db.QueryRowContext(ctx, "SELECT role FROM users WHERE id = $1", id).Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrUserNotFound

		}

		return "", fmt.Errorf("get user %s role: %w", id, err)
	}

	return domain.Role(role), nil
}
