package tasks

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

type fakeRepository struct {
	createErr  error
	getListErr error
	deleteErr  error

	createCalls  int
	getListCalls int
	deleteCalls  int

	gotID          uuid.UUID
	gotAuthorID    uuid.UUID
	gotAssigneeID  uuid.UUID
	gotName        string
	gotDescription string
	gotPriority    string
	gotStatus      *string

	gotListAuthorID   uuid.UUID
	gotListAssigneeID *uuid.UUID
	gotListStatus     domain.TaskStatus
	gotListPriority   domain.TaskPriority
	gotListSearch     string
	gotListPage       int
	gotListLimit      int

	taskList  []domain.Task
	taskTotal int

	gotDeleteID       uuid.UUID
	gotDeleteAuthorID uuid.UUID
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
	ctx context.Context,
	authorID uuid.UUID,
	assigneeID *uuid.UUID,
	status domain.TaskStatus,
	priority domain.TaskPriority,
	search string,
	page, limit int,
) ([]domain.Task, int, error) {
	f.getListCalls++

	f.gotListAuthorID = authorID
	f.gotListAssigneeID = assigneeID
	f.gotListStatus = status
	f.gotListPriority = priority
	f.gotListSearch = search
	f.gotListPage = page
	f.gotListLimit = limit

	if f.getListErr != nil {
		return []domain.Task{}, 0, f.getListErr
	}

	return f.taskList, f.taskTotal, nil
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
