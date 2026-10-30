package tasks

import (
	"context"

	"github.com/google/uuid"
)

func (uc *tasksUseCases) DeleteTask(ctx context.Context, id, authorID uuid.UUID) error {
	return uc.repo.DeleteTask(ctx, id, authorID)
}
