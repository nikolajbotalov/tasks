package tasks

import (
	"TaskFlow/internal/domain"
	"TaskFlow/internal/testutils"
	"errors"
	"testing"

	"github.com/google/uuid"
)

var (
	testAuthorID   = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	testAssigneeID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
)

func TestCreateTask(t *testing.T) {
	t.Parallel()

	ptr := func(s string) *string { return &s }

	repoErr := errors.New("connection refused")

	tests := []struct {
		name string

		assigneeID  uuid.UUID
		status      string
		priority    string
		taskName    string
		description string
		repoErr     error

		wantErr      error
		wantAssignee uuid.UUID
		wantStatus   *string
		wantPriority string
	}{
		{
			name:         "defaults for empty status and priority",
			assigneeID:   uuid.Nil,
			taskName:     "task",
			description:  "description",
			wantAssignee: testAuthorID,
			wantStatus:   ptr(string(domain.ToDo)),
			wantPriority: string(domain.TaskLowPriority),
		},
		{
			name:         "keeps explicit assignee",
			assigneeID:   testAssigneeID,
			taskName:     "task",
			description:  "description",
			wantAssignee: testAssigneeID,
			wantStatus:   ptr(string(domain.ToDo)),
			wantPriority: string(domain.TaskLowPriority),
		},
		{
			name:         "keeps explicit status and priority",
			assigneeID:   testAssigneeID,
			status:       string(domain.InReview),
			priority:     string(domain.TaskHighPriority),
			taskName:     "task",
			description:  "description",
			wantAssignee: testAssigneeID,
			wantStatus:   ptr(string(domain.InReview)),
			wantPriority: string(domain.TaskHighPriority),
		},
		{
			name:         "repo error is wrapped and still matchable",
			assigneeID:   testAssigneeID,
			taskName:     "task",
			description:  "description",
			repoErr:      repoErr,
			wantErr:      repoErr,
			wantAssignee: testAssigneeID,
			wantStatus:   ptr(string(domain.ToDo)),
			wantPriority: string(domain.TaskLowPriority),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{createErr: tt.repoErr}
			uc := NewTasksUseCases(repo)

			got, err := uc.CreateTask(
				t.Context(),
				testAuthorID,
				tt.assigneeID,
				tt.taskName,
				tt.description,
				tt.status,
				tt.priority,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if repo.createCalls != 1 {
				t.Fatalf("repo.CreateTask called %d times, want 1", repo.createCalls)
			}

			if repo.gotID == uuid.Nil {
				t.Errorf("generated ID is nil")
			}

			if repo.gotAuthorID != testAuthorID {
				t.Errorf("authorID = %s, want %s", repo.gotAuthorID, testAuthorID)
			}

			if repo.gotAssigneeID != tt.wantAssignee {
				t.Errorf("assigneeID = %s, want %s", repo.gotAssigneeID, tt.wantAssignee)
			}

			if repo.gotName != tt.taskName {
				t.Errorf("name = %q, want %q", repo.gotName, tt.taskName)
			}

			if repo.gotDescription != tt.description {
				t.Errorf("description = %q, want %q", repo.gotDescription, tt.description)
			}

			if repo.gotPriority != tt.wantPriority {
				t.Errorf("priority = %q, want %q", repo.gotPriority, tt.wantPriority)
			}

			testutils.AssertOptionalString(t, "status", repo.gotStatus, tt.wantStatus)

			if tt.wantErr != nil {
				return
			}

			if got.ID != repo.gotID {
				t.Errorf("task.ID = %s, want %s", got.ID, repo.gotID)
			}
		})
	}
}
