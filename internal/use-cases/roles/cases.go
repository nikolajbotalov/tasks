package roles

import (
	"TaskFlow/internal/domain"

	"github.com/google/uuid"
)

type roleRepository interface {
	UpdateUserRole(userID uuid.UUID, role domain.Role) error
}

type UseCases interface {
	UpdateUserRole(userID uuid.UUID, role string) error
}

type rolesUseCases struct {
	repo roleRepository
}

func NewRolesUseCases(repo roleRepository) UseCases {
	return &rolesUseCases{
		repo: repo,
	}
}
