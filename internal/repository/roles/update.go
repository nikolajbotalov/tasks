package roles

import (
	"TaskFlow/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) UpdateUserRole(userID uuid.UUID, role domain.Role) error {
	res, err := r.db.Exec("UPDATE users SET role = $1, updated_at = $2 WHERE id = $3", role, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("update user %s with role %s: %w", userID, role, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update user %s with role %s, affected rows: %w", userID, role, err)
	}

	if affected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
