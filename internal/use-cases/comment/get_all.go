package comment

import (
	"TaskFlow/internal/domain"

	"github.com/google/uuid"
)

func (uc *commentUseCases) GetCommentList(taskID uuid.UUID, page, limit int) (domain.GetTaskCommentsResponse, error) {
	list, total, err := uc.repo.GetCommentList(taskID, page, limit)
	if err != nil {
		return domain.GetTaskCommentsResponse{}, err
	}

	return domain.GetTaskCommentsResponse{
		Comments: list,
		Page:     page,
		Limit:    limit,
		Total:    total,
	}, nil
}
