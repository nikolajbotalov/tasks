package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *authUseCases) CreateUser(ctx context.Context, email, password string) (uuid.UUID, error) {
	normalizedEmail := domain.NormalizeEmail(email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %w", domain.ErrHashedPassword, err)
	}

	id := uuid.New()

	err = uc.repo.CreateUser(ctx, id.String(), normalizedEmail, string(hashedPassword))
	if err != nil {
		if errors.Is(err, domain.ErrEmailExists) {
			return uuid.Nil, fmt.Errorf("%w: %w", domain.ErrEmailExists, err)
		}

		return uuid.Nil, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}
