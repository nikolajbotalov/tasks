package roles

import (
	"TaskFlow/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) UpdateUserRole(ctx context.Context, userID uuid.UUID, role domain.Role) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin role update for user %s: %w", userID, err)
	}
	defer func() { _ = tx.Rollback() }()

	adminIDs, err := lockAdmin(ctx, tx)
	if err != nil {
		return err
	}

	if len(adminIDs) == 1 && adminIDs[0] == userID && role != domain.RoleAdmin {
		return domain.ErrLastAdmin
	}

	res, err := tx.ExecContext(ctx,
		"UPDATE users SET role = $1, updated_at = $2 WHERE id = $3",
		role, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("update user %s role: %w", userID, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update user %s role %s, affected rows: %w", userID, role, err)
	}
	if affected == 0 {
		return domain.ErrUserNotFound
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit role update for user %s: %w", userID, err)
	}

	return nil
}

func lockAdmin(ctx context.Context, tx *sql.Tx) ([]uuid.UUID, error) {
	rows, err := tx.QueryContext(ctx,
		"SELECT id FROM users WHERE role = $1 ORDER BY id FOR NO KEY UPDATE",
		domain.RoleAdmin)
	if err != nil {
		return nil, fmt.Errorf("lock admins: %w", err)
	}
	defer rows.Close()

	ids := make([]uuid.UUID, 0)
	for rows.Next() {
		var id uuid.UUID
		if err = rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan admin id: %w", err)
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate admin ids: %w", err)
	}

	return ids, nil
}
