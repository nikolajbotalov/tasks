package tasks

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

var testTaskID = uuid.MustParse("11111111-1111-1111-1111-111111111111")

func TestDeleteTask(t *testing.T) {
	t.Parallel()

	deleteErr := errors.New("connection refused")

	tests := []struct {
		name    string
		repoErr error
		wantErr error
	}{
		{
			name: "forwards id and authorID to repository",
		},
		{
			name:    "repo error is propagated",
			repoErr: deleteErr,
			wantErr: deleteErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{deleteErr: tt.repoErr}
			uc := NewTasksUseCases(repo)

			err := uc.DeleteTask(t.Context(), testTaskID, testAuthorID)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if repo.deleteCalls != 1 {
				t.Fatalf("repo.DeleteTask called %d times, want 1", repo.deleteCalls)
			}

			if repo.gotDeleteID != testTaskID {
				t.Errorf("id = %s, want %s", repo.gotDeleteID, testTaskID)
			}

			if repo.gotDeleteAuthorID != testAuthorID {
				t.Errorf("authorID = %s, want %s", repo.gotDeleteAuthorID, testAuthorID)
			}
		})
	}
}
