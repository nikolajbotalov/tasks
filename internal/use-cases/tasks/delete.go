package tasks

import (
	"github.com/google/uuid"
)

func (uc *tasksUseCases) DeleteTask(id, authorID uuid.UUID) error {
	return uc.repo.DeleteTask(id, authorID)
}
