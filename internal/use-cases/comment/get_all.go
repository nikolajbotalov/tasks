package comment

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

func (uc *commentUseCases) GetCommentList(
	ctx context.Context,
	taskID uuid.UUID,
	page, limit int,
) (domain.GetTaskCommentsResponse, error) {
	list, total, err := uc.repo.GetCommentList(ctx, taskID, page, limit)
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
