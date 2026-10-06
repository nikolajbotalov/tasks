package tasks

import (
	"github.com/google/uuid"
)

func (uc *tasksUseCases) DeleteTask(id, authorID uuid.UUID) error {
	err := uc.repo.DeleteTask(id, authorID)
	if err != nil {
		return err
	}

	return nil
}
