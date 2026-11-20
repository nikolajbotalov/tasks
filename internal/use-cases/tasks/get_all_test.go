package tasks

import (
	"TaskFlow/internal/domain"
	"TaskFlow/internal/testutils"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestGetAllTasks(t *testing.T) {
	t.Parallel()

	assigneeID := uuid.MustParse("55555555-5555-5555-5555-555555555555")
	ptrUUID := func(u uuid.UUID) *uuid.UUID { return &u }

	repoErr := errors.New("connection refused")

	repoTasks := []domain.Task{
		{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), Name: "first", AuthorID: testAuthorID},
		{ID: uuid.MustParse("22222222-2222-2222-2222-222222222222"), Name: "second", AuthorID: testAuthorID},
	}

	tests := []struct {
		name string

		assigneeID *uuid.UUID
		status     domain.TaskStatus
		priority   domain.TaskPriority
		search     string
		page       int
		limit      int

		tasks   []domain.Task
		total   int
		repoErr error

		wantErr error
	}{
		{
			name:  "return tasks with pagination",
			page:  3,
			limit: 25,
			tasks: repoTasks,
			total: 42,
		},
		{
			name:  "empty list",
			page:  1,
			limit: 10,
			tasks: []domain.Task{},
			total: 0,
		},
		{
			name:       "forward all filters",
			assigneeID: ptrUUID(assigneeID),
			status:     domain.ToDo,
			priority:   domain.TaskHighPriority,
			search:     "first",
			page:       2,
			limit:      10,
			tasks:      repoTasks,
			total:      2,
		},
		{
			name:    "repo error is wrapped",
			page:    1,
			limit:   10,
			repoErr: repoErr,
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{
				getListErr: tt.repoErr,
				taskList:   tt.tasks,
				taskTotal:  tt.total,
			}
			uc := NewTasksUseCases(repo)

			got, err := uc.GetAllTasks(
				t.Context(),
				testAuthorID,
				tt.assigneeID,
				tt.status,
				tt.priority,
				tt.search,
				tt.page,
				tt.limit,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if repo.getListCalls != 1 {
				t.Fatalf("repo.GetTaskList called %d times, want 1", repo.getListCalls)
			}

			if repo.gotListAuthorID != testAuthorID {
				t.Errorf("authorID = %s, want %s", repo.gotListAuthorID, testAuthorID)
			}

			testutils.AssertOptionalUUID(t, "assigneeID", repo.gotListAssigneeID, tt.assigneeID)

			if repo.gotListStatus != tt.status {
				t.Errorf("status = %s, want %s", repo.gotListStatus, tt.status)
			}

			if repo.gotListPriority != tt.priority {
				t.Errorf("priority = %s, want %s", repo.gotListPriority, tt.priority)
			}

			if repo.gotListSearch != tt.search {
				t.Errorf("search = %s, want %s", repo.gotListSearch, tt.search)
			}

			if repo.gotListPage != tt.page {
				t.Errorf("page = %d, want %d", repo.gotListPage, tt.page)
			}

			if repo.gotListLimit != tt.limit {
				t.Errorf("limit = %d, want %d", repo.gotListLimit, tt.limit)
			}

			if tt.wantErr != nil {
				return
			}

			if got.Page != tt.page {
				t.Errorf("response.Page = %d, want %d", got.Page, tt.page)
			}

			if got.Limit != tt.limit {
				t.Errorf("response.Limit = %d, want %d", got.Limit, tt.limit)
			}

			if got.Total != tt.total {
				t.Errorf("response.Total = %d, want %d", got.Total, tt.total)
			}

			if !reflect.DeepEqual(got.Tasks, tt.tasks) {
				t.Errorf("response.Tasks = %+v, want %+v", got.Tasks, tt.tasks)
			}
		})
	}
}
