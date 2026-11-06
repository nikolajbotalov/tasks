package auth

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, id, email, password string) error
	FindUserByEmail(ctx context.Context, email string) (domain.UserRow, error)
	GetAuthState(ctx context.Context, id uuid.UUID) (domain.AuthState, error)
	CreateRefreshSession(ctx context.Context, session domain.RefreshSession) error
	GetRefreshSessionState(ctx context.Context, id uuid.UUID) (domain.RefreshSessionState, error)
	RotateRefreshSession(
		ctx context.Context,
		oldID, newID, userID, familyID uuid.UUID,
		expiresAt time.Time,
	) (bool, domain.RefreshSessionState, error)
	RevokeRefreshFamilyAndBumpVersion(ctx context.Context, userID, familyID uuid.UUID) error
	RevokeAllSessionsAndBumpVersion(ctx context.Context, userID uuid.UUID) error
}

type UseCases interface {
	CreateUser(ctx context.Context, email, password string) (uuid.UUID, error)
	LoginUser(ctx context.Context, email, password string) (domain.TokenPair, uuid.UUID, error)
	RefreshToken(ctx context.Context, refreshToken string) (domain.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	GetAuthState(ctx context.Context, id uuid.UUID) (domain.AuthState, error)
}

type authUseCases struct {
	repo   Repository
	config *config.Config
}

func NewAuthUseCases(repo Repository, config *config.Config) UseCases {
	return &authUseCases{
		repo:   repo,
		config: config,
	}
}
