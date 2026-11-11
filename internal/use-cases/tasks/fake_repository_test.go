package tasks

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

type fakeRepository struct {
	createErr error

	createCalls int

	gotID          uuid.UUID
	gotAuthorID    uuid.UUID
	gotAssigneeID  uuid.UUID
	gotName        string
	gotDescription string
	gotPriority    string
	gotStatus      *string
}

var _ Repository = (*fakeRepository)(nil)

func (f *fakeRepository) CreateTask(
	ctx context.Context,
	id, authorID, assigneeID uuid.UUID,
	name, description, priority string,
	status *string,
) (domain.Task, error) {
	f.createCalls++

	f.gotID = id
	f.gotAuthorID = authorID
	f.gotAssigneeID = assigneeID
	f.gotName = name
	f.gotDescription = description
	f.gotPriority = priority
	f.gotStatus = status

	if f.createErr != nil {
		return domain.Task{}, f.createErr
	}

	return domain.Task{
		ID:          id,
		Name:        name,
		Description: description,
		AuthorID:    authorID,
		AssigneeID:  assigneeID,
		Status:      status,
		Priority:    &priority,
	}, nil
}

func (f *fakeRepository) GetTaskList(
	context.Context, uuid.UUID, *uuid.UUID, domain.TaskStatus, domain.TaskPriority, string, int, int,
) ([]domain.Task, int, error) {
	panic("unexpected call: GetTaskList")
}

func (f *fakeRepository) GetTaskByID(context.Context, uuid.UUID) (domain.Task, error) {
	panic("unexpected call: GetTaskByID")
}

func (f *fakeRepository) UpdateTask(
	context.Context, uuid.UUID, *uuid.UUID, *string, *string, *string, *string,
) error {
	panic("unexpected call: UpdateTask")
}

func (f *fakeRepository) DeleteTask(context.Context, uuid.UUID, uuid.UUID) error {
	panic("unexpected call: DeleteTask")
}
