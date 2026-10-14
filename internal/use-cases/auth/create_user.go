package auth

import (
	"TaskFlow/internal/domain"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (uc *authUseCases) CreateUser(id, username, email, password string) error {
	trimEmail := strings.TrimSpace(email)
	lowerEmail := strings.ToLower(trimEmail)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%v: %w", domain.ErrHashedPassword, err)
	}

	err = uc.repo.CreateUser(id, username, lowerEmail, string(hashedPassword))
	if err != nil {
		if errors.Is(err, domain.ErrEmailExists) {
			return fmt.Errorf("%v: %w", domain.ErrEmailExists, err)
		}

		return fmt.Errorf("create user: %w", err)
	}

	return nil
}
