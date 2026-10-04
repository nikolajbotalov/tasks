package auth

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/domain"
)

type Repository interface {
	CreateUser(id, username, email, password string) error
	FindUserByEmail(email string) (domain.UserRow, error)
}

type UseCases interface {
	CreateUser(id, username, email, password string) error
	LoginUser(email, password string) (domain.TokenPair, error)
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
