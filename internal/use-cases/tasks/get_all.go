package tasks

import (
	"TaskFlow/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

func (uc *tasksUseCases) GetAllTasks(
	authorID uuid.UUID,
	assigneeID *uuid.UUID,
	status domain.TaskStatus,
	priority domain.TaskPriority,
	page, limit int,
) (domain.GetAllTasksResponse, error) {
	userTasks, total, err := uc.repo.GetTaskList(authorID, assigneeID, status, priority, page, limit)
	if err != nil {
		return domain.GetAllTasksResponse{}, fmt.Errorf("get task list: %w", err)
	}

	return domain.GetAllTasksResponse{
		Tasks: userTasks,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}
