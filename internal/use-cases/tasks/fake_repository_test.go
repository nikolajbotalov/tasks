package tasks

import (
	"TaskFlow/internal/domain"
	"context"

	"github.com/google/uuid"
)

type fakeRepository struct {
	createErr  error
	getListErr error
	getByIDErr error
	deleteErr  error
	updateErr  error

	createCalls  int
	getListCalls int
	getByIDCalls int
	deleteCalls  int
	updateCalls  int

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

	gotUpdateID          uuid.UUID
	gotUpdateAssigneeID  *uuid.UUID
	gotUpdateName        *string
	gotUpdateDescription *string
	gotUpdateStatus      *string
	gotUpdatePriority    *string

	gotGetByID uuid.UUID
	taskByID   domain.Task
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

func (f *fakeRepository) GetTaskByID(cxt context.Context, id uuid.UUID) (domain.Task, error) {
	f.getByIDCalls++
	f.gotGetByID = id

	return f.taskByID, f.getByIDErr
}

func (f *fakeRepository) UpdateTask(
	ctx context.Context,
	id uuid.UUID,
	assigneeID *uuid.UUID,
	name, description, status, priority *string,
) error {
	f.updateCalls++

	f.gotUpdateID = id
	f.gotUpdateAssigneeID = assigneeID
	f.gotUpdateName = name
	f.gotUpdateDescription = description
	f.gotUpdateStatus = status
	f.gotUpdatePriority = priority

	return f.updateErr
}

func (f *fakeRepository) DeleteTask(ctx context.Context, id, authorID uuid.UUID) error {
	f.deleteCalls++

	f.gotDeleteID = id
	f.gotDeleteAuthorID = authorID

	return f.deleteErr
}
