package roles

import "TaskFlow/internal/domain"

func (uc *rolesUseCases) GetAllUsersRoles() (domain.GetUserRolesResponse, error) {
	roles, err := uc.repo.GetAllUsersRoles()
	if err != nil {
		return domain.GetUserRolesResponse{}, err
	}

	return domain.GetUserRolesResponse{
		Roles: roles,
	}, nil
}
