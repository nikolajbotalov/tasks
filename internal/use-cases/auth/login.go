package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *authUseCases) LoginUser(ctx context.Context, email, password string) (domain.TokenPair, uuid.UUID, error) {
	normalizedEmail := domain.NormalizeEmail(email)

	user, err := uc.repo.FindUserByEmail(ctx, normalizedEmail)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.TokenPair{}, uuid.Nil, domain.ErrInvalidCredentials
		}

		return domain.TokenPair{}, uuid.Nil, fmt.Errorf("login: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return domain.TokenPair{}, uuid.Nil, domain.ErrInvalidCredentials
		}

		return domain.TokenPair{}, uuid.Nil, fmt.Errorf("compare password: %w", err)
	}

	tokens, err := uc.generateTokenPair(user.ID, user.TokenVersion)
	if err != nil {
		return domain.TokenPair{}, uuid.Nil, fmt.Errorf("generate token pair: %w", err)
	}

	return tokens, user.ID, nil
}
