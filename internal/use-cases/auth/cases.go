package auth

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(id, username, email, password string) error
	FindUserByEmail(email string) (domain.UserRow, error)
	GetUserRole(ctx context.Context, id uuid.UUID) (domain.Role, error)
}

type UseCases interface {
	CreateUser(id, username, email, password string) error
	LoginUser(email, password string) (domain.TokenPair, error)
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
