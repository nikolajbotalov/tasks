package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *authRepository) GetRefreshSessionState(ctx context.Context, id uuid.UUID) (domain.RefreshSessionState, error) {
	state, err := scanRefreshSession(r.db.QueryRowContext(ctx,
		`SELECT user_id, family_id, used_at, revoked_at, expires_at
				FROM refresh_sessions
				WHERE id = $1`, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.RefreshSessionState{}, domain.ErrInvalidRefreshToken
		}

		return domain.RefreshSessionState{}, fmt.Errorf("get refresh session %s: %w", id, err)
	}

	return state, nil
}
