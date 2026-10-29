package profiles

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

func (uc *profilesUseCases) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (domain.Profile, error) {
	return uc.repo.GetProfileByUserID(ctx, userID)
}
