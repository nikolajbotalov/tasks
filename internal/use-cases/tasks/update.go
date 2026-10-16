package tasks

import (
	"github.com/google/uuid"
)

func (uc *tasksUseCases) UpdateTask(id, assigneeID uuid.UUID, name, description, status, priority *string) error {
	if err := uc.repo.UpdateTask(id, assigneeID, name, description, status, priority); err != nil {
		return err
	}
	return nil
}
