package auth

import (
	"TaskFlow/internal/domain"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (uc *authUseCases) CreateUser(id, username, email, password string) error {
	normalizedEmail := domain.NormalizeEmail(email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrHashedPassword, err)
	}

	err = uc.repo.CreateUser(id, username, normalizedEmail, string(hashedPassword))
	if err != nil {
		if errors.Is(err, domain.ErrEmailExists) {
			return fmt.Errorf("%w: %w", domain.ErrEmailExists, err)
		}

		return fmt.Errorf("create user: %w", err)
	}

	return nil
}
