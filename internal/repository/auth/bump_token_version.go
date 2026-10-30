package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *authRepository) BumpTokenVersion(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx,
		"UPDATE users SET token_version = token_version + 1, updated_at = now() WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("bump token version for user %s: %w", id, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("bump token version for user %s, affected rows: %w", id, err)
	}
	if affected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
