package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"fmt"
)

func (r *authRepository) CreateRefreshSession(ctx context.Context, session domain.RefreshSession) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO refresh_sessions (id, user_id, family_id, expires_at)
				VALUES ($1, $2, $3, $4)`,
		session.ID, session.UserID, session.FamilyID, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("create refresh session %s: %w", session.ID, err)
	}

	return nil
}
