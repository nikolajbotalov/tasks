package comment

import (
	"TaskFlow/internal/domain"

	"github.com/google/uuid"
)

func (uc *commentUseCases) CreateComment(taskID, authorID uuid.UUID, text string) (domain.TaskComment, error) {
	commentID := uuid.New()

	comment, err := uc.repo.CreateComment(commentID, taskID, authorID, text)
	if err != nil {
		return domain.TaskComment{}, err
	}

	return comment, nil
}
