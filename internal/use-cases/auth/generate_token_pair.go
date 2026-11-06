package auth

import (
	"TaskFlow/internal/domain"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (uc *authUseCases) generateTokenPair(userID uuid.UUID, version int) (domain.TokenPair, domain.RefreshSession, error) {
	return uc.generateTokenPairInFamily(userID, uuid.New(), version)
}

func (uc *authUseCases) generateTokenPairInFamily(
	userID, familyID uuid.UUID,
	version int,
) (domain.TokenPair, domain.RefreshSession, error) {
	jti := uuid.New()
	expiresAt := time.Now().Add(uc.config.JWT.RefreshTokenTTL)

	accessToken, err := uc.signToken(
		userID,
		uuid.Nil,
		uuid.Nil,
		version,
		domain.JWTTypeAccess,
		uc.config.JWT.AccessTokenTTL,
	)
	if err != nil {
		return domain.TokenPair{}, domain.RefreshSession{}, err
	}

	refreshToken, err := uc.signToken(
		userID,
		jti,
		familyID,
		version,
		domain.JWTTypeRefresh,
		uc.config.JWT.RefreshTokenTTL,
	)
	if err != nil {
		return domain.TokenPair{}, domain.RefreshSession{}, err
	}

	return domain.TokenPair{
			AccessToken:      accessToken,
			RefreshToken:     refreshToken,
			RefreshExpiresAt: expiresAt,
		},
		domain.RefreshSession{
			ID:        jti,
			UserID:    userID,
			FamilyID:  familyID,
			ExpiresAt: expiresAt,
		}, nil
}

func (uc *authUseCases) signToken(
	userID, jti, familyID uuid.UUID,
	version int,
	typ string,
	ttl time.Duration,
) (string, error) {
	claims := jwt.MapClaims{
		domain.JWTID:           userID,
		domain.JWTTokenVersion: version,
		domain.JWTType:         typ,
		"exp":                  jwt.NewNumericDate(time.Now().Add(ttl)),
	}

	if jti != uuid.Nil {
		claims[domain.JWTJTI] = jti.String()
	}
	if familyID != uuid.Nil {
		claims[domain.JWTFAM] = familyID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(uc.config.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("sign %s token: %w", typ, err)
	}

	return signed, nil
}
