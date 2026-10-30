package roles

import (
	"TaskFlow/internal/domain"
	"context"
)

func (uc *rolesUseCases) GetAllUsersRoles(ctx context.Context) (domain.GetUserRolesResponse, error) {
	roles, err := uc.repo.GetAllUsersRoles(ctx)
	if err != nil {
		return domain.GetUserRolesResponse{}, err
	}

	return domain.GetUserRolesResponse{
		Roles: roles,
	}, nil
}
