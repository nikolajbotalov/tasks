package auth

import (
	"TaskFlow/internal/config"
	"TaskFlow/internal/domain"
	"errors"
	"testing"

	"github.com/google/uuid"
)

var testUserID = uuid.MustParse("11111111-1111-1111-1111-111111111111")

func TestGetAuthState(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("connection refused")

	validState := domain.AuthState{Role: domain.RoleAdmin, TokenVersion: 3}

	tests := []struct {
		name      string
		repoState domain.AuthState
		repoErr   error
		wantErr   error
		wantState domain.AuthState
	}{
		{
			name:      "returns valid state",
			repoState: validState,
			wantState: validState,
		},
		{
			name:      "rejects unknown role from db",
			repoState: domain.AuthState{Role: "owner", TokenVersion: 1},
			wantErr:   domain.ErrUnknownRole,
		},
		{
			name:      "rejects empty role from db",
			repoState: domain.AuthState{Role: "", TokenVersion: 1},
			wantErr:   domain.ErrUnknownRole,
		},
		{
			name:    "repo error is propagated",
			repoErr: repoErr,
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{
				authState: tt.repoState,
				authErr:   tt.repoErr,
			}
			uc := NewAuthUseCases(repo, &config.Config{})

			got, err := uc.GetAuthState(t.Context(), testUserID)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if repo.authCalls != 1 {
				t.Fatalf("repo.GetAuthState called %d times, want 1", repo.authCalls)
			}

			if repo.gotAuthID != testUserID {
				t.Errorf("id = %s, want %s", repo.gotAuthID, testUserID)
			}

			if tt.wantErr != nil {
				if got != (domain.AuthState{}) {
					t.Errorf("got = %+v, want zero state on error", got)
				}
				return
			}

			if got != tt.wantState {
				t.Errorf("got = %+v, want %+v", got, tt.wantState)
			}
		})
	}
}
