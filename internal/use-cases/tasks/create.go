package tasks

import (
	"TaskFlow/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

func (uc *tasksUseCases) CreateTask(authorID uuid.UUID, name, description string) (domain.Task, error) {
	newUUID := uuid.New()

	task, err := uc.repo.CreateTask(newUUID, authorID, name, description)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}
