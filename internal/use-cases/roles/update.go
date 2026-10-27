package roles

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

func (uc *rolesUseCases) UpdateUserRole(ctx context.Context, targetID, curUserID uuid.UUID, role string) error {
	r := domain.Role(role)

	if !r.IsValid() {
		return domain.ErrUnknownRole
	}

	if targetID == curUserID {
		return domain.ErrSelfRoleChange
	}

	return uc.repo.UpdateUserRole(ctx, targetID, r)
}
