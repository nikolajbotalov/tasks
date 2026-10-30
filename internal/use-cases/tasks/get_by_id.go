package tasks

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

func (uc *tasksUseCases) GetTaskByID(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	return uc.repo.GetTaskByID(ctx, id)
}
