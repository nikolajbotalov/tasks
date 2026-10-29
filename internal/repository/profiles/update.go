package profiles

import (
	"TaskFlow/internal/domain"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *profilesRepository) UpdateProfile(
	ctx context.Context,
	username, firstName, lastName *string,
	userID uuid.UUID,
) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE profiles 
				SET username = COALESCE($1, username),
					first_name = COALESCE($2, first_name),
					last_name = COALESCE($3, last_name),
					updated_at = $4
				WHERE user_id = $5`, username, firstName, lastName, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("update profile: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update profile with userID %s, affected rows: %w", userID, err)
	}

	if affected == 0 {
		return domain.ErrProfileNotFound
	}

	return nil
}
