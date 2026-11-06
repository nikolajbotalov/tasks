package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"errors"
	"fmt"
	"time"
)

func (uc *authUseCases) Logout(ctx context.Context, refreshToken string) error {
	claims, err := uc.parseRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	state, err := uc.repo.GetAuthState(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("get auth state for user %s: %w", claims.UserID, err)
	}

	if claims.Version != state.TokenVersion {
		return domain.ErrInvalidRefreshToken
	}

	session, err := uc.repo.GetRefreshSessionState(ctx, claims.JTI)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidRefreshToken) {
			return domain.ErrInvalidRefreshToken
		}
		return fmt.Errorf("get refresh session %s: %w", claims.JTI, err)
	}

	if session.UserID != claims.UserID ||
		session.UsedAt != nil ||
		session.RevokedAt != nil ||
		time.Now().After(session.ExpiresAt) {
		return domain.ErrInvalidRefreshToken
	}

	if err = uc.repo.RevokeAllSessionsAndBumpVersion(ctx, claims.UserID); err != nil {
		return fmt.Errorf("revoke sessions for user %s: %w", claims.UserID, err)
	}

	return nil
}
