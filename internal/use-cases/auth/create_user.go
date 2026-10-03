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
		fmt.Printf("failed to hash password, email: %v\n", lowerEmail)
		return domain.ErrHashedPassword
	}

	err = uc.repo.CreateUser(id, username, lowerEmail, string(hashedPassword))
	if err != nil {
		if errors.Is(err, domain.ErrEmailExists) {
			fmt.Printf("user %v with email %v, cannot creating\n", username, lowerEmail)
			return domain.ErrEmailExists
		}

		fmt.Printf("failed to create user %v, with email: %v\n", username, err)
		return domain.ErrInternal
	}

	return nil
}
