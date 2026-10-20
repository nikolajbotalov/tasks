package tasks

import (
	"TaskFlow/internal/domain"

	"github.com/google/uuid"
)

func (uc *tasksUseCases) GetTaskByID(id uuid.UUID) (domain.Task, error) {
	return uc.repo.GetTaskByID(id)
}
