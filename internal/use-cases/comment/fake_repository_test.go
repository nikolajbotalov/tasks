package comment

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

type fakeRepository struct {
	createErr  error
	getListErr error

	createCalls  int
	getListCalls int

	gotID       uuid.UUID
	gotTaskID   uuid.UUID
	gotAuthorID uuid.UUID
	gotText     string

	gotListTaskID uuid.UUID
	gotListPage   int
	gotListLimit  int

	commentList  []domain.TaskComment
	commentTotal int
}

var _ Repository = (*fakeRepository)(nil)

func (f *fakeRepository) CreateComment(
	ctx context.Context,
	ID, taskID, authorID uuid.UUID,
	text string,
) (domain.TaskComment, error) {
	f.createCalls++

	f.gotID = ID
	f.gotTaskID = taskID
	f.gotAuthorID = authorID
	f.gotText = text

	if f.createErr != nil {
		return domain.TaskComment{}, f.createErr
	}

	return domain.TaskComment{
		ID:         ID,
		TaskID:     taskID,
		AuthorID:   authorID,
		AuthorName: "user",
		Text:       text,
	}, nil
}

func (f *fakeRepository) GetCommentList(
	ctx context.Context,
	taskID uuid.UUID,
	page, limit int,
) ([]domain.TaskComment, int, error) {
	f.getListCalls++

	f.gotListTaskID = taskID
	f.gotListPage = page
	f.gotListLimit = limit

	if f.getListErr != nil {
		return []domain.TaskComment{}, 0, f.getListErr
	}

	return f.commentList, f.commentTotal, nil
}
