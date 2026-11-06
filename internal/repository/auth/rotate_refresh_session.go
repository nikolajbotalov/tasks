package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *authRepository) RotateRefreshSession(
	ctx context.Context,
	oldID, newID, userID, familyID uuid.UUID,
	expiresAt time.Time,
) (bool, domain.RefreshSessionState, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, domain.RefreshSessionState{}, fmt.Errorf("begin rotate refresh session: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var updatedID uuid.UUID
	err = tx.QueryRowContext(ctx,
		`UPDATE refresh_sessions
				SET used_at = now()
				WHERE id = $1
				  AND user_id = $2
				  AND family_id = $3
				  AND used_at IS NULL
				  AND revoked_at IS NULL
				  AND expires_at > now()
				RETURNING id`,
		oldID, userID, familyID,
	).Scan(&updatedID)

	switch {
	case err == nil:
		if _, err = tx.ExecContext(ctx,
			`INSERT INTO refresh_sessions (id, user_id, family_id, expires_at)
					VALUES ($1, $2, $3, $4)`,
			newID, userID, familyID, expiresAt,
		); err != nil {
			return false, domain.RefreshSessionState{}, fmt.Errorf("insert rotated refresh session: %w", err)
		}

		if err = tx.Commit(); err != nil {
			return false, domain.RefreshSessionState{}, fmt.Errorf("commit rotate refresh session: %w", err)
		}

		return true, domain.RefreshSessionState{}, nil

	case errors.Is(err, sql.ErrNoRows):
		state, stateErr := scanRefreshSession(tx.QueryRowContext(ctx,
			`SELECT user_id, family_id, used_at, revoked_at, expires_at
					FROM refresh_sessions
					WHERE id = $1`, oldID))
		if stateErr != nil {
			if errors.Is(stateErr, sql.ErrNoRows) {
				return false, domain.RefreshSessionState{}, nil
			}

			return false, domain.RefreshSessionState{}, fmt.Errorf("get refresh session %s: %w", oldID, stateErr)
		}

		return false, state, nil

	default:
		return false, domain.RefreshSessionState{}, fmt.Errorf("mark refresh session %s used: %w", oldID, err)
	}
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRefreshSession(row rowScanner) (domain.RefreshSessionState, error) {
	var (
		state     domain.RefreshSessionState
		usedAt    sql.NullTime
		revokedAt sql.NullTime
	)

	if err := row.Scan(&state.UserID, &state.FamilyID, &usedAt, &revokedAt, &state.ExpiresAt); err != nil {
		return domain.RefreshSessionState{}, err
	}

	if usedAt.Valid {
		t := usedAt.Time
		state.UsedAt = &t
	}
	if revokedAt.Valid {
		t := revokedAt.Time
		state.RevokedAt = &t
	}

	return state, nil
}
