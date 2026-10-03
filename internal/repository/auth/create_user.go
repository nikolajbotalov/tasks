package auth

import (
	"TaskFlow/internal/domain"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (r *authRepository) CreateUser(id, username, email, password string) error {
	curTime := time.Now()

	err := r.db.QueryRow("INSERT INTO users (id, email, username, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (email) DO NOTHING RETURNING id",
		id, email, username, password, curTime, curTime).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("user with email %v already exists\n", email)
			return domain.ErrEmailExists
		}

		fmt.Printf("failed to insert user with email %s by error: %v\n", email, err)
		return domain.ErrInternal
	}

	return nil
}
