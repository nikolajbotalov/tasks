package domain

import (
	"time"

	"github.com/google/uuid"
)

type TokenPair struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"-"`
	RefreshExpiresAt time.Time `json:"-"`
}

type RefreshSession struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FamilyID  uuid.UUID
	ExpiresAt time.Time
}

type RefreshSessionState struct {
	UserID    uuid.UUID
	FamilyID  uuid.UUID
	UsedAt    *time.Time
	RevokedAt *time.Time
	ExpiresAt time.Time
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterResponse struct {
	ID string `json:"id"`
}

type AuthState struct {
	Role         Role
	TokenVersion int
}
