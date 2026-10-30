package auth

import (
	"TaskFlow/internal/domain"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (uc *authUseCases) generateTokenPair(userID uuid.UUID, version int) (domain.TokenPair, error) {
	accessToken, err := uc.signToken(userID, version, domain.JWTTypeAccess, uc.config.JWT.AccessTokenTTL)
	if err != nil {
		return domain.TokenPair{}, err
	}

	refreshToken, err := uc.signToken(userID, version, domain.JWTTypeRefresh, uc.config.JWT.RefreshTokenTTL)
	if err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *authUseCases) signToken(userID uuid.UUID, version int, typ string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		domain.JWTID:           userID,
		domain.JWTTokenVersion: version,
		domain.JWTType:         typ,
		"exp":                  jwt.NewNumericDate(time.Now().Add(ttl)),
	})

	signed, err := token.SignedString([]byte(uc.config.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("sign %s token: %w", typ, err)
	}

	return signed, nil
}
