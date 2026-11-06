package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const refreshReuseGrace = 10 * time.Second

type refreshTokenClaims struct {
	UserID   uuid.UUID
	Version  int
	JTI      uuid.UUID
	FamilyID uuid.UUID
}

func (uc *authUseCases) RefreshToken(ctx context.Context, refreshToken string) (domain.TokenPair, error) {
	claims, err := uc.parseRefreshToken(refreshToken)
	if err != nil {
		return domain.TokenPair{}, err
	}

	state, err := uc.repo.GetAuthState(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.TokenPair{}, domain.ErrInvalidRefreshToken
		}

		return domain.TokenPair{}, fmt.Errorf("get auth state for user %s: %w", claims.UserID, err)
	}
	if claims.Version != state.TokenVersion {
		return domain.TokenPair{}, domain.ErrInvalidRefreshToken
	}

	tokens, session, err := uc.generateTokenPairInFamily(claims.UserID, claims.FamilyID, state.TokenVersion)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate token pair: %w", err)
	}

	rotated, oldSession, err := uc.repo.RotateRefreshSession(
		ctx,
		claims.JTI,
		session.ID,
		claims.UserID,
		claims.FamilyID,
		session.ExpiresAt,
	)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("rotate refresh session %s: %w", claims.JTI, err)
	}
	if rotated {
		return tokens, nil
	}

	if oldSession.UserID != claims.UserID || oldSession.FamilyID != claims.FamilyID {
		return domain.TokenPair{}, domain.ErrInvalidRefreshToken
	}
	if oldSession.RevokedAt != nil || time.Now().After(oldSession.ExpiresAt) {
		return domain.TokenPair{}, domain.ErrInvalidRefreshToken
	}
	if oldSession.UsedAt == nil {
		return domain.TokenPair{}, domain.ErrInvalidRefreshToken
	}

	if time.Since(*oldSession.UsedAt) <= refreshReuseGrace {
		return domain.TokenPair{}, domain.ErrRefreshInProgress
	}

	if err = uc.repo.RevokeRefreshFamilyAndBumpVersion(ctx, claims.UserID, claims.FamilyID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.TokenPair{}, domain.ErrInvalidRefreshToken
		}
		return domain.TokenPair{}, fmt.Errorf("revoke refresh family %s: %w", claims.FamilyID, err)
	}

	return domain.TokenPair{}, domain.ErrInvalidRefreshToken
}

func (uc *authUseCases) parseRefreshToken(tokenString string) (refreshTokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(uc.config.JWT.Secret), nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil || !token.Valid {
		return refreshTokenClaims{}, domain.ErrInvalidRefreshToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return refreshTokenClaims{}, domain.ErrInvalidRefreshToken
	}

	if typ, ok := claims[domain.JWTType].(string); !ok || typ != domain.JWTTypeRefresh {
		return refreshTokenClaims{}, domain.ErrInvalidRefreshToken
	}

	id, ok := claims[domain.JWTID].(string)
	if !ok {
		return refreshTokenClaims{}, domain.ErrInvalidRefreshToken
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return refreshTokenClaims{}, domain.ErrInvalidRefreshToken
	}

	jti, ok := claims[domain.JWTJTI].(string)
	if !ok {
		return refreshTokenClaims{}, domain.ErrInvalidRefreshToken
	}
	sessionID, err := uuid.Parse(jti)
	if err != nil {
		return refreshTokenClaims{}, domain.ErrInvalidRefreshToken
	}

	family, ok := claims[domain.JWTFAM].(string)
	if !ok {
		return refreshTokenClaims{}, domain.ErrInvalidRefreshToken
	}
	familyID, err := uuid.Parse(family)
	if err != nil {
		return refreshTokenClaims{}, domain.ErrInvalidRefreshToken
	}

	version, ok := claims[domain.JWTTokenVersion].(float64)
	if !ok {
		return refreshTokenClaims{}, domain.ErrInvalidRefreshToken
	}

	return refreshTokenClaims{
		UserID:   userID,
		Version:  int(version),
		JTI:      sessionID,
		FamilyID: familyID,
	}, nil
}
