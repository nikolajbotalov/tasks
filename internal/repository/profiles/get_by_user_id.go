package profiles

import (
	"TaskFlow/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *profilesRepository) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (domain.Profile, error) {
	var profile domain.Profile

	err := r.db.QueryRowContext(ctx,
		"SELECT user_id, username, first_name, last_name, created_at, updated_at FROM profiles WHERE user_id = $1",
		userID).Scan(&profile.UserID, &profile.Username, &profile.FirstName, &profile.LastName, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Profile{}, domain.ErrUserNotFound
		}

		return domain.Profile{}, fmt.Errorf("get profile with user id %s: %w", userID, err)
	}

	return profile, nil
}
