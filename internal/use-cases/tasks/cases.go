package tasks

import (
	"TaskFlow/internal/domain"

	"github.com/google/uuid"
)

type Repository interface {
	GetTaskList(
		authorID uuid.UUID,
		assigneeID *uuid.UUID,
		status domain.TaskStatus,
		priority domain.TaskPriority,
		search string,
		page, limit int,
	) ([]domain.Task, int, error)
	GetTaskByID(id uuid.UUID) (domain.Task, error)
	CreateTask(id, authorID, assigneeID uuid.UUID, name, description, priority string, status *string) (domain.Task, error)
	UpdateTask(id uuid.UUID, assigneeID *uuid.UUID, name, description, status, priority *string) error
	DeleteTask(id, authorID uuid.UUID) error
}

type UseCases interface {
	GetAllTasks(
		authorID uuid.UUID,
		assigneeID *uuid.UUID,
		status domain.TaskStatus,
		priority domain.TaskPriority,
		search string,
		page, limit int,
	) (domain.GetAllTasksResponse, error)
	GetTaskByID(id uuid.UUID) (domain.Task, error)
	CreateTask(authorID, assigneeID uuid.UUID, name, description, status, priority string) (domain.Task, error)
	UpdateTask(id uuid.UUID, assigneeID *uuid.UUID, name, description, status, priority *string) error
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
