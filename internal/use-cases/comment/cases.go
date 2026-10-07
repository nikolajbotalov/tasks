package comment

import (
	"TaskFlow/internal/domain"

	"github.com/google/uuid"
)

type Repository interface {
	GetCommentList(taskID uuid.UUID, page, limit int) ([]domain.TaskComment, int, error)
	CreateComment(ID, taskID, authorID uuid.UUID, text string) (domain.TaskComment, error)
}

type UseCases interface {
	GetCommentList(taskID uuid.UUID, page, limit int) (domain.GetTaskCommentsResponse, error)
	CreateComment(taskID, authorID uuid.UUID, text string) (domain.TaskComment, error)
}

type commentUseCases struct {
	repo Repository
}

func NewCommentUseCases(repo Repository) UseCases {
	return &commentUseCases{
		repo: repo,
	}
}
