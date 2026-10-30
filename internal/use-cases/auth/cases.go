package auth

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, id, email, password string) error
	FindUserByEmail(ctx context.Context, email string) (domain.UserRow, error)
	GetAuthState(ctx context.Context, id uuid.UUID) (domain.AuthState, error)
	BumpTokenVersion(ctx context.Context, id uuid.UUID) error
}

type UseCases interface {
	CreateUser(ctx context.Context, email, password string) (uuid.UUID, error)
	LoginUser(ctx context.Context, email, password string) (domain.TokenPair, uuid.UUID, error)
	RefreshToken(ctx context.Context, refreshToken string) (domain.TokenPair, error)
	Logout(ctx context.Context, userID uuid.UUID) error
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
