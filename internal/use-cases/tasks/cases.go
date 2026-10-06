package tasks

import (
	"TaskFlow/internal/domain"

	"github.com/google/uuid"
)

type Repository interface {
	GetTaskList(authorID uuid.UUID, page, limit int) ([]domain.Task, int, error)
	CreateTask(id, authorID uuid.UUID, name, description string) (domain.Task, error)
	DeleteTask(id, authorID uuid.UUID) error
}

type UseCases interface {
	GetAllTasks(authorID uuid.UUID, page, limit int) (domain.GetAllTasksResponse, error)
	CreateTask(authorID uuid.UUID, name, description string) (domain.Task, error)
	DeleteTask(id, authorID uuid.UUID) error
}

type tasksUseCases struct {
	repo Repository
}

func NewTasksUseCases(repo Repository) UseCases {
	return &tasksUseCases{
		repo: repo,
	}
}
