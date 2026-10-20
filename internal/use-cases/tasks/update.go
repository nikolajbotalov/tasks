package tasks

import (
	"github.com/google/uuid"
)

func (uc *tasksUseCases) UpdateTask(
	id uuid.UUID,
	assigneeID *uuid.UUID,
	name, description, status, priority *string,
) error {
	return uc.repo.UpdateTask(id, assigneeID, name, description, status, priority)
}
