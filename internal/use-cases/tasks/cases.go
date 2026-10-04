package tasks

import (
	"TaskFlow/internal/domain"

	"github.com/google/uuid"
)

type Repository interface {
	CreateTask(id, authorID uuid.UUID, name, description string) (domain.Task, error)
}

type UseCases interface {
	CreateTask(authorID uuid.UUID, name, description string) (domain.Task, error)
}

type tasksUseCases struct {
	repo Repository
}

func NewTasksUseCases(repo Repository) UseCases {
	return &tasksUseCases{
		repo: repo,
	}
}
