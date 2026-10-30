package tasks

import (
	"TaskFlow/internal/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (uc *tasksUseCases) GetAllTasks(
	ctx context.Context,
	authorID uuid.UUID,
	assigneeID *uuid.UUID,
	status domain.TaskStatus,
	priority domain.TaskPriority,
	search string,
	page, limit int,
) (domain.GetAllTasksResponse, error) {
	userTasks, total, err := uc.repo.GetTaskList(ctx, authorID, assigneeID, status, priority, search, page, limit)
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
