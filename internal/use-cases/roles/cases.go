package roles

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

type roleRepository interface {
	GetAllUsersRoles() ([]domain.Role, error)
	UpdateUserRole(ctx context.Context, userID uuid.UUID, role domain.Role) error
}

type UseCases interface {
	GetAllUsersRoles() (domain.GetUserRolesResponse, error)
	UpdateUserRole(ctx context.Context, targetID, curUserID uuid.UUID, role string) error
}

type rolesUseCases struct {
	repo roleRepository
}

func NewRolesUseCases(repo roleRepository) UseCases {
	return &rolesUseCases{
		repo: repo,
	}
}
