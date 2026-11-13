package roles

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

type fakeRepository struct {
	updateErr error

	updateCalls int

	gotTargetID uuid.UUID
	gotRole     domain.Role
}

var _ roleRepository = (*fakeRepository)(nil)

func (f *fakeRepository) GetAllUsersRoles(context.Context) ([]domain.Role, error) {
	panic("unexpected call: GetAllUsersRoles")
}

func (f *fakeRepository) UpdateUserRole(
	ctx context.Context,
	userID uuid.UUID,
	role domain.Role,
) error {
	f.updateCalls++

	f.gotTargetID = userID
	f.gotRole = role

	return f.updateErr
}
