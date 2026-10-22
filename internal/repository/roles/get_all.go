package roles

import (
	"TaskFlow/internal/domain"
	"fmt"
)

func (r *Repository) GetAllUsersRoles() ([]domain.Role, error) {
	roles := make([]domain.Role, 0)

	rows, err := r.db.Query("SELECT name FROM roles ORDER BY name")
	if err != nil {
		return []domain.Role{}, fmt.Errorf("get all user roles: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var role domain.Role
		err = rows.Scan(&role)
		if err != nil {
			return []domain.Role{}, fmt.Errorf("scan user role row: %w", err)
		}
		roles = append(roles, role)
	}

	err = rows.Err()
	if err != nil {
		return []domain.Role{}, fmt.Errorf("iterate user role rows: %w", err)
	}

	return roles, nil
}
