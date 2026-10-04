package auth

import (
	"TaskFlow/internal/domain"
	"database/sql"
	"errors"
)

func (r *authRepository) FindUserByEmail(email string) (domain.UserRow, error) {
	var userRow domain.UserRow

	rowErr := r.db.QueryRow("SELECT id, password FROM users WHERE email = $1", email).
		Scan(&userRow.ID, &userRow.Password)
	if rowErr != nil {
		if errors.Is(rowErr, sql.ErrNoRows) {
			return userRow, domain.ErrUserNotFound
		}

		return userRow, domain.ErrInternal
	}

	return userRow, nil
}
