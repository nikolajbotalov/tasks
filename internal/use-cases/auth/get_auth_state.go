package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (uc *authUseCases) GetAuthState(ctx context.Context, id uuid.UUID) (domain.AuthState, error) {
	state, err := uc.repo.GetAuthState(ctx, id)
	if err != nil {
		return domain.AuthState{}, fmt.Errorf("get auth state: %w", err)
	}

	if !state.Role.IsValid() {
		return domain.AuthState{}, fmt.Errorf("%w: %s", domain.ErrUnknownRole, state.Role)
	}

	return state, nil
}
