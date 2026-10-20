package tasks

import (
	"TaskFlow/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

func (uc *tasksUseCases) CreateTask(
	authorID, assigneeID uuid.UUID,
	name, description, status, priority string,
) (domain.Task, error) {
	newUUID := uuid.New()

	assignee := assigneeID
	if assigneeID == uuid.Nil {
		assignee = authorID
	}

	taskStatus := status
	if taskStatus == "" {
		taskStatus = string(domain.ToDo)
	}

	taskPriority := priority
	if taskPriority == "" {
		taskPriority = string(domain.TaskLowPriority)
	}

	task, err := uc.repo.CreateTask(newUUID, authorID, assignee, name, description, taskPriority, &taskStatus)
	if err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}
