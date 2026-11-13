package comment

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

type fakeRepository struct {
	createErr error

	createCalls int

	gotID       uuid.UUID
	gotTaskID   uuid.UUID
	gotAuthorID uuid.UUID
	gotText     string
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
	context.Context, uuid.UUID, int, int,
) ([]domain.TaskComment, int, error) {
	panic("unexpected call: GetCommentList")
}
