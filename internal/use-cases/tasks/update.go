package tasks

import (
	"fmt"

	"github.com/google/uuid"
)

func (uc *tasksUseCases) UpdateTask(id uuid.UUID, name, description *string) error {
	if err := uc.repo.UpdateTask(id, name, description); err != nil {
		return fmt.Errorf("update task %s: %w", id, err)
	}
	return nil
}
