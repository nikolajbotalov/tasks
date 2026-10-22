package roles

import (
	"TaskFlow/internal/domain"

	"github.com/google/uuid"
)

func (uc *rolesUseCases) UpdateUserRole(userID uuid.UUID, role string) error {
	r := domain.Role(role)

	if !r.IsValid() {
		return domain.ErrUnknownRole
	}

	return uc.repo.UpdateUserRole(userID, r)
}
