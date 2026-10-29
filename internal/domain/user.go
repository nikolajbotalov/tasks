package domain

import (
	"github.com/google/uuid"
)

type UserRow struct {
	ID       uuid.UUID
	Password string
}

type LoginUserResponse struct {
	ID           uuid.UUID `json:"id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}
