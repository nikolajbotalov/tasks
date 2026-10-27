package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (uc *authUseCases) GetCurrentRole(ctx context.Context, id uuid.UUID) (domain.Role, error) {
	role, err := uc.repo.GetUserRole(ctx, id)
	if err != nil {
		return "", fmt.Errorf("get current role: %w", err)
	}

	if !role.IsValid() {
		return "", fmt.Errorf("%w: %s", domain.ErrUnknownRole, role)
	}

	return role, nil
}
