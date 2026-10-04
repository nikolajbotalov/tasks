package auth

import (
	"database/sql"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *authRepository {
	return &authRepository{
		db: db,
	}
}
