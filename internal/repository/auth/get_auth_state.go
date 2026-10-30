package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *authRepository) GetAuthState(ctx context.Context, id uuid.UUID) (domain.AuthState, error) {
	var (
		state domain.AuthState
		role  string
	)

	err := r.db.QueryRowContext(ctx, "SELECT role, token_version FROM users WHERE id = $1", id).
		Scan(&role, &state.TokenVersion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.AuthState{}, domain.ErrUserNotFound
		}

		return domain.AuthState{}, fmt.Errorf("get user %s auth state: %w", id, err)
	}

	state.Role = domain.Role(role)

	return state, nil
}
