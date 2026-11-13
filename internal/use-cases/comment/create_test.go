package comment

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

var (
	testTaskID   = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	testAuthorID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
)

func TestCreateComment(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("connection refused")

	tests := []struct {
		name    string
		text    string
		repoErr error
		wantErr error
	}{
		{
			name: "forwards generated id and arguments",
			text: "comment",
		},
		{
			name: "keeps multiline text",
			text: "line one\nline two",
		},
		{
			name:    "repo error is propagated",
			text:    "comment",
			repoErr: repoErr,
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{createErr: tt.repoErr}
			uc := NewCommentUseCases(repo)

			got, err := uc.CreateComment(t.Context(), testTaskID, testAuthorID, tt.text)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want = %v", err, tt.wantErr)
			}

			if repo.createCalls != 1 {
				t.Fatalf("repo.CreateComment called %d times, want 1", repo.createCalls)
			}

			if repo.gotID == uuid.Nil {
				t.Error("generated ID is nil")
			}

			if repo.gotID == testTaskID || repo.gotID == testAuthorID {
				t.Errorf("generated ID collides with an arguments: %s", repo.gotID)
			}

			if repo.gotTaskID != testTaskID {
				t.Errorf("taskID = %s, want %s", repo.gotTaskID, testTaskID)
			}

			if repo.gotAuthorID != testAuthorID {
				t.Errorf("authorID = %s, want %s", repo.gotAuthorID, testAuthorID)
			}

			if repo.gotText != tt.text {
				t.Errorf("text = %s, want %s", repo.gotText, tt.text)
			}

			if tt.wantErr != nil {
				return
			}

			if got.ID != repo.gotID {
				t.Errorf("comment.ID = %s, want %s", got.ID, repo.gotID)
			}

			if got.AuthorID != testAuthorID {
				t.Errorf("comment.AuthorID = %s, want %s", got.AuthorID, testAuthorID)
			}

			if got.Text != tt.text {
				t.Errorf("comment.Text = %q, want %q", got.Text, tt.text)
			}
		})
	}
}
