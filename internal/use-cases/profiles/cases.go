package profiles

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (domain.Profile, error)
	UpdateProfile(ctx context.Context, username, firstName, lastName *string, userID uuid.UUID) error
}

type UseCases interface {
	GetProfileByUserID(ctx context.Context, userID uuid.UUID) (domain.Profile, error)
	UpdateProfile(ctx context.Context, username, firstName, lastName *string, userID, callerID uuid.UUID) error
}

type profilesUseCases struct {
	repo Repository
}

func NewProfilesUseCases(repo Repository) UseCases {
	return &profilesUseCases{repo: repo}
}
