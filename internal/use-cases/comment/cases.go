package comment

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetCommentList(ctx context.Context, taskID uuid.UUID, page, limit int) ([]domain.TaskComment, int, error)
	CreateComment(ctx context.Context, ID, taskID, authorID uuid.UUID, text string) (domain.TaskComment, error)
}

type UseCases interface {
	GetCommentList(ctx context.Context, taskID uuid.UUID, page, limit int) (domain.GetTaskCommentsResponse, error)
	CreateComment(ctx context.Context, taskID, authorID uuid.UUID, text string) (domain.TaskComment, error)
}

type commentUseCases struct {
	repo Repository
}

func NewCommentUseCases(repo Repository) UseCases {
	return &commentUseCases{
		repo: repo,
	}
}
