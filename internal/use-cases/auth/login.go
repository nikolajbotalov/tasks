package auth

import (
	"TaskFlow/internal/domain"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (uc *authUseCases) LoginUser(email, password string) (domain.TokenPair, error) {
	normalizedEmail := domain.NormalizeEmail(email)

	user, err := uc.repo.FindUserByEmail(normalizedEmail)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.TokenPair{}, domain.ErrInvalidCredentials
		}

		return domain.TokenPair{}, fmt.Errorf("login: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return domain.TokenPair{}, domain.ErrInvalidCredentials
		}

		return domain.TokenPair{}, fmt.Errorf("compare password: %w", err)
	}

	tokens, err := uc.generateTokenPair(user.ID, user.Role)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate token pair: %w", err)
	}

	return tokens, nil
}
