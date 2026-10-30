package comment

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

func (uc *commentUseCases) CreateComment(
	ctx context.Context,
	taskID, authorID uuid.UUID,
	text string,
) (domain.TaskComment, error) {
	commentID := uuid.New()

	comment, err := uc.repo.CreateComment(ctx, commentID, taskID, authorID, text)
	if err != nil {
		return domain.TaskComment{}, err
	}

	return comment, nil
}
