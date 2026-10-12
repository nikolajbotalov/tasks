package permissions

import (
	"TaskFlow/internal/domain"
	"fmt"
)

func (r *permissionsRepository) GetAll() (domain.RolePermissions, error) {
	rows, err := r.db.Query("SELECT role, permission FROM role_permissions")
	if err != nil {
		return nil, fmt.Errorf("get role permissions: %w", err)
	}
	defer rows.Close()

	result := make(domain.RolePermissions)

	for rows.Next() {
		var role, permission string
		if err := rows.Scan(&role, &permission); err != nil {
			return nil, fmt.Errorf("scan role permissions: %w", err)
		}

		dbRole := domain.Role(role)
		if result[dbRole] == nil {
			result[dbRole] = make(map[domain.Permission]struct{})
		}
		result[dbRole][domain.Permission(permission)] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate role permissions: %w", err)
	}

	return result, nil
}
