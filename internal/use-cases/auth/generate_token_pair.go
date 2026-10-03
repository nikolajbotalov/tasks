package auth

import (
	"TaskFlow/internal/domain"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (uc *authUseCases) generateTokenPair(userID uuid.UUID) (domain.TokenPair, error) {
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": jwt.NewNumericDate(time.Now().Add(uc.config.JWT.AccessTokenTTL)),
	})

	secret := []byte(uc.config.JWT.Secret)
	accessToken, err := access.SignedString(secret)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	refresh := make([]byte, 32)
	if _, err = rand.Read(refresh); err != nil {
		return domain.TokenPair{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	refreshToken := base64.URLEncoding.EncodeToString(refresh)

	return domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
