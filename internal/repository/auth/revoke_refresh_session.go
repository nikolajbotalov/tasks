package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (r *authRepository) RevokeRefreshFamilyAndBumpVersion(ctx context.Context, userID, familyID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin revoke refresh family: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err = tx.ExecContext(ctx,
		`UPDATE refresh_sessions
				SET revoked_at = now()
				WHERE user_id = $1 AND family_id = $2 AND revoked_at IS NULL`,
		userID, familyID,
	); err != nil {
		return fmt.Errorf("revoke refresh family %s: %w", familyID, err)
	}

	if err = bumpTokenVersionTx(ctx, tx, userID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit revoke refresh family: %w", err)
	}

	return nil
}

func (r *authRepository) RevokeAllSessionsAndBumpVersion(ctx context.Context, userID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin revoke all sessions: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err = tx.ExecContext(ctx,
		`UPDATE refresh_sessions
				SET revoked_at = now()
				WHERE user_id = $1 AND revoked_at IS NULL`,
		userID,
	); err != nil {
		return fmt.Errorf("revoke sessions for user %s: %w", userID, err)
	}

	if err = bumpTokenVersionTx(ctx, tx, userID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit revoke all sessions: %w", err)
	}

	return nil
}

func bumpTokenVersionTx(ctx context.Context, tx *sql.Tx, userID uuid.UUID) error {
	res, err := tx.ExecContext(ctx,
		`UPDATE users
				SET token_version = token_version + 1, updated_at = now()
				WHERE id = $1`, userID)
	if err != nil {
		return fmt.Errorf("bump token version for user %s: %w", userID, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("bump token version for user %s, affected rows: %w", userID, err)
	}
	if affected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
