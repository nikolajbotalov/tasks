package tasks

import (
	"TaskFlow/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

func (uc *tasksUseCases) GetAllTasks(authorID uuid.UUID, page, limit int) (domain.GetAllTasksResponse, error) {
	userTasks, total, err := uc.repo.GetTaskList(authorID, page, limit)
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
