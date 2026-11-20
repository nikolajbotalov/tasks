package tasks

import (
	"TaskFlow/internal/domain"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

var testGetTaskID = uuid.MustParse("11111111-1111-1111-1111-111111111111")

func TestGetTaskByID(t *testing.T) {
	t.Parallel()

	repoTask := domain.Task{
		ID:          testGetTaskID,
		Name:        "task",
		Description: "description",
		AuthorID:    testAuthorID,
		AssigneeID:  testAssigneeID,
	}

	tests := []struct {
		name     string
		repoTask domain.Task
		repoErr  error
		wantErr  error
	}{
		{
			name:     "return task from repository",
			repoTask: repoTask,
		},
		{
			name:    "task not found is propagated",
			repoErr: domain.ErrTaskNotFound,
			wantErr: domain.ErrTaskNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{
				taskByID:   tt.repoTask,
				getByIDErr: tt.repoErr,
			}
			uc := NewTasksUseCases(repo)

			got, err := uc.GetTaskByID(t.Context(), testGetTaskID)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want = %v", err, tt.wantErr)
			}

			if repo.getByIDCalls != 1 {
				t.Fatalf("repo.GetTaskByID called %d times, want 1", repo.getByIDCalls)
			}

			if repo.gotGetByID != testGetTaskID {
				t.Errorf("id = %s, want %s", repo.gotGetByID, testGetTaskID)
			}

			if tt.wantErr != nil {
				return
			}

			if !reflect.DeepEqual(got, tt.repoTask) {
				t.Errorf("task = %+v, want %+v", got, tt.repoTask)
			}
		})
	}
}
