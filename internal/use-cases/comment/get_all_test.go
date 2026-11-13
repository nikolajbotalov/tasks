package comment

import (
	"TaskFlow/internal/domain"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestGetCommentList(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("connection refused")

	repoComments := []domain.TaskComment{
		{
			ID:         uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			TaskID:     testTaskID,
			AuthorID:   testAuthorID,
			AuthorName: "user",
			Text:       "first comment",
		},
		{
			ID:         uuid.MustParse("22222222-2222-2222-2222-222222222222"),
			TaskID:     testTaskID,
			AuthorID:   testAuthorID,
			AuthorName: "user",
			Text:       "second comment",
		},
	}

	tests := []struct {
		name string

		page  int
		limit int

		repoComments []domain.TaskComment
		repoTotal    int
		repoErr      error

		wantErr error
	}{
		{
			name:         "return comments with pagination",
			page:         2,
			limit:        3,
			repoComments: repoComments,
			repoTotal:    7,
		},
		{
			name:         "empty list is not an error",
			page:         1,
			limit:        10,
			repoComments: []domain.TaskComment{},
			repoTotal:    0,
		},
		{
			name:    "repo error is propagated",
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
				getListErr:   tt.repoErr,
				commentList:  tt.repoComments,
				commentTotal: tt.repoTotal,
			}
			uc := NewCommentUseCases(repo)

			got, err := uc.GetCommentList(t.Context(), testTaskID, tt.page, tt.limit)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if repo.getListCalls != 1 {
				t.Fatalf("repo.GetCommentList called %d times, want 1", repo.getListCalls)
			}

			if repo.gotListTaskID != testTaskID {
				t.Errorf("taskID = %s, want %s", repo.gotListTaskID, testTaskID)
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

			if got.Total != tt.repoTotal {
				t.Errorf("response.Total = %d, want %d", got.Total, tt.repoTotal)
			}

			if !reflect.DeepEqual(got.Comments, tt.repoComments) {
				t.Errorf("response.Comments = %+v, want %+v", got.Comments, tt.repoComments)
			}
		})
	}
}
