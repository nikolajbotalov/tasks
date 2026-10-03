package auth

import (
	"TaskFlow/internal/domain"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (uc *authUseCases) LoginUser(email, password string) (domain.TokenPair, error) {
	trimEmail := strings.TrimSpace(email)
	lowerEmail := strings.ToLower(trimEmail)

	user, err := uc.repo.FindUserByEmail(lowerEmail)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			fmt.Printf("user with email %s not found by error: %v\n", email, err)
			return domain.TokenPair{}, domain.ErrInvalidCredentials
		}

		fmt.Printf("failed to find user with email %s by error: %v\n", email, err)
		return domain.TokenPair{}, domain.ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			fmt.Println(err)
			return domain.TokenPair{}, domain.ErrInvalidCredentials
		}

		fmt.Println(err)
		return domain.TokenPair{}, domain.ErrInternal
	}

	tokens, err := uc.generateTokenPair(user.ID)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("failed to generate token pair: %w", err)
	}

	return tokens, nil
}
