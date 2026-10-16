package tasks

import (
	"TaskFlow/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

func (uc *tasksUseCases) CreateTask(authorID, assigneeID uuid.UUID, name, description, status, priority string) (domain.Task, error) {
	newUUID := uuid.New()

	var (
		assignee     uuid.UUID
		taskStatus   string
		taskPriority string
	)

	if assigneeID == uuid.Nil {
		assignee = authorID
	}

	if status == "" {
		taskStatus = string(domain.ToDo)
	}

	if priority == "" {
		taskPriority = string(domain.TaskLowPriority)
	}

	task, err := uc.repo.CreateTask(newUUID, assignee, authorID, name, description, taskPriority, &taskStatus)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}
