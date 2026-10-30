package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (uc *authUseCases) RefreshToken(ctx context.Context, refreshToken string) (domain.TokenPair, error) {
	userID, version, err := uc.parseRefreshToken(refreshToken)
	if err != nil {
		return domain.TokenPair{}, err
	}

	state, err := uc.repo.GetAuthState(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.TokenPair{}, domain.ErrInvalidRefreshToken
		}

		return domain.TokenPair{}, fmt.Errorf("get auth state for user %s: %w", userID, err)
	}

	if version != state.TokenVersion {
		return domain.TokenPair{}, domain.ErrInvalidRefreshToken
	}

	tokens, err := uc.generateTokenPair(userID, state.TokenVersion)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate token pair: %w", err)
	}

	return tokens, nil
}

func (uc *authUseCases) parseRefreshToken(refreshToken string) (uuid.UUID, int, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(uc.config.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, 0, domain.ErrInvalidRefreshToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, 0, domain.ErrInvalidRefreshToken
	}

	typ, ok := claims[domain.JWTType].(string)
	if !ok || typ != domain.JWTTypeRefresh {
		return uuid.Nil, 0, domain.ErrInvalidRefreshToken
	}

	id, ok := claims[domain.JWTID].(string)
	if !ok {
		return uuid.Nil, 0, domain.ErrInvalidRefreshToken
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, 0, domain.ErrInvalidRefreshToken
	}

	version, ok := claims[domain.JWTTokenVersion].(float64)
	if !ok {
		return uuid.Nil, 0, domain.ErrInvalidRefreshToken
	}

	return userID, int(version), nil
}
