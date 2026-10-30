package tasks

import (
	"context"

	"github.com/google/uuid"
)

func (uc *tasksUseCases) UpdateTask(
	ctx context.Context,
	id uuid.UUID,
	assigneeID *uuid.UUID,
	name, description, status, priority *string,
) error {
	return uc.repo.UpdateTask(ctx, id, assigneeID, name, description, status, priority)
}
