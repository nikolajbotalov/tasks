package auth

import (
	"TaskFlow/internal/domain"
	"database/sql"
	"errors"
	"fmt"
)

func (r *authRepository) FindUserByEmail(email string) (UserRow, error) {
	var userRow UserRow

	rowErr := r.db.QueryRow("SELECT id, password FROM users WHERE email = $1", email).
		Scan(&userRow.ID, &userRow.Password)
	if rowErr != nil {
		if errors.Is(rowErr, sql.ErrNoRows) {
			fmt.Printf("user not found: %v\n", rowErr)
			return userRow, domain.ErrUserNotFound
		}

		fmt.Printf("internal server error: %v\n", rowErr)
		return userRow, domain.ErrInternal
	}

	return userRow, nil
}
