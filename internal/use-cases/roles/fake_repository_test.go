package roles

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

type fakeRepository struct {
	updateErr error
	getAllErr error

	updateCalls int
	getAllCalls int

	gotTargetID uuid.UUID
	gotRole     domain.Role

	roleList []domain.Role
}

var _ roleRepository = (*fakeRepository)(nil)

func (f *fakeRepository) GetAllUsersRoles(ctx context.Context) ([]domain.Role, error) {
	f.getAllCalls++

	if f.getAllErr != nil {
		return []domain.Role{}, f.getAllErr
	}

	return f.roleList, nil
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
