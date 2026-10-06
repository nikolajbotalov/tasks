package tasks

import (
	"TaskFlow/internal/domain"

	"github.com/google/uuid"
)

func (uc *tasksUseCases) GetTaskByID(id uuid.UUID) (domain.Task, error) {
	task, err := uc.repo.GetTaskByID(id)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}
