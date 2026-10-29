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
	GetUserRole(ctx context.Context, id uuid.UUID) (domain.Role, error)
}

type UseCases interface {
	CreateUser(ctx context.Context, email, password string) (uuid.UUID, error)
	LoginUser(ctx context.Context, email, password string) (domain.TokenPair, uuid.UUID, error)
	GetCurrentRole(ctx context.Context, id uuid.UUID) (domain.Role, error)
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
