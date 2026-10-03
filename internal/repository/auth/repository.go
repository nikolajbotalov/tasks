package auth

import (
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(id, username, email, password string) error
	FindUserByEmail(email string) (UserRow, error)
}

type UserRow struct {
	ID       uuid.UUID `json:"id"`
	Password string    `json:"password"`
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) Repository {
	return &authRepository{
		db: db,
	}
}
