package tasks

import (
	"TaskFlow/internal/domain"
	"TaskFlow/internal/testutils"
	"errors"
	"testing"

	"github.com/google/uuid"
)

var testUpdateTaskID = uuid.MustParse("11111111-1111-1111-1111-111111111111")

func TestUpdateTask(t *testing.T) {
	t.Parallel()

	ptr := func(s string) *string { return &s }
	ptrUUID := func(u uuid.UUID) *uuid.UUID { return &u }

	repoErr := errors.New("connection refused")

	tests := []struct {
		name string

		assigneeID  *uuid.UUID
		taskName    *string
		description *string
		status      *string
		priority    *string
		repoErr     error

		wantErr error
	}{
		{
			name:        "forwards all fields",
			assigneeID:  ptrUUID(testAssigneeID),
			taskName:    ptr("test"),
			description: ptr("description"),
			status:      ptr(string(domain.InProgress)),
			priority:    ptr(string(domain.TaskLowPriority)),
		},
		{
			name:     "partial update keeps omitted fields nil",
			taskName: ptr("test"),
		},
		{
			name:     "explicit empty string is preserved",
			taskName: ptr(""),
			status:   ptr(""),
		},
		{
			name:     "repo error is propagated",
			taskName: ptr("test"),
			repoErr:  repoErr,
			wantErr:  repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{updateErr: tt.repoErr}
			uc := NewTasksUseCases(repo)

			err := uc.UpdateTask(
				t.Context(),
				testUpdateTaskID,
				tt.assigneeID,
				tt.taskName,
				tt.description,
				tt.status,
				tt.priority,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if repo.updateCalls != 1 {
				t.Fatalf("repo.UpdateTask called %d times, want 1", repo.updateCalls)
			}

			if repo.gotUpdateID != testUpdateTaskID {
				t.Errorf("id = %s, want %s", repo.gotUpdateID, testUpdateTaskID)
			}

			testutils.AssertOptionalUUID(t, "assigneeID", repo.gotUpdateAssigneeID, tt.assigneeID)
			testutils.AssertOptionalString(t, "name", repo.gotUpdateName, tt.taskName)
			testutils.AssertOptionalString(t, "description", repo.gotUpdateDescription, tt.description)
			testutils.AssertOptionalString(t, "status", repo.gotUpdateStatus, tt.status)
			testutils.AssertOptionalString(t, "priority", repo.gotUpdatePriority, tt.priority)
		})
	}
}
