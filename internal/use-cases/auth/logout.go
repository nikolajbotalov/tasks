package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (uc *authUseCases) Logout(ctx context.Context, userID uuid.UUID) error {
	if err := uc.repo.BumpTokenVersion(ctx, userID); err != nil {
		return fmt.Errorf("logout user %s: %w", userID, err)
	}

	return nil
}
