package profiles

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

type fakeRepository struct {
	updateErr error

	getCalls    int
	updateCalls int

	gotUsername  *string
	gotFirstName *string
	gotLastName  *string
	gotUserID    uuid.UUID
}

var _ Repository = (*fakeRepository)(nil)

func (f *fakeRepository) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (domain.Profile, error) {
	f.getCalls++
	return domain.Profile{}, nil
}

func (f *fakeRepository) UpdateProfile(
	ctx context.Context,
	username, firstName, lastName *string,
	userID uuid.UUID,
) error {
	f.updateCalls++
	f.gotUsername, f.gotFirstName, f.gotLastName, f.gotUserID = username, firstName, lastName, userID

	return f.updateErr
}
