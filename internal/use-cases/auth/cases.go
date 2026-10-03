package auth

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/domain"
	"TaskFlow/internal/repository/auth"
)

type UseCases interface {
	CreateUser(id, username, email, password string) error
	LoginUser(email, password string) (domain.TokenPair, error)
}

type authUseCases struct {
	repo   auth.Repository
	config *config.Config
}

func NewAuthUseCases(repo auth.Repository, config *config.Config) UseCases {
	return &authUseCases{
		repo:   repo,
		config: config,
	}
}
