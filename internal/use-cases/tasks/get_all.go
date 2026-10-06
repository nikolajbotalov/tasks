package tasks

import (
	"TaskFlow/internal/domain"

	"github.com/google/uuid"
)

func (uc *tasksUseCases) GetAllTasks(authorID uuid.UUID, page, limit int) (domain.GetAllTasksResponse, error) {
	userTasks, total, err := uc.repo.GetTaskList(authorID, page, limit)
	if err != nil {
		return domain.GetAllTasksResponse{}, err
	}

	return domain.GetAllTasksResponse{
		Tasks: userTasks,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}
