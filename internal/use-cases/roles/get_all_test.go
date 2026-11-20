package roles

import (
	"TaskFlow/internal/domain"
	"errors"
	"reflect"
	"testing"
)

func TestGetAllUsersRoles(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("connection refused")

	repoRoles := []domain.Role{domain.RoleAdmin, domain.RoleUser, domain.RoleViewer}

	tests := []struct {
		name      string
		repoRoles []domain.Role
		repoErr   error
		wantErr   error
	}{
		{
			name:      "returns roles from repository",
			repoRoles: repoRoles,
		},
		{
			name:      "empty list is not an error",
			repoRoles: []domain.Role{},
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
				roleList:  tt.repoRoles,
				getAllErr: tt.repoErr,
			}
			uc := NewRolesUseCases(repo)

			got, err := uc.GetAllUsersRoles(t.Context())

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if repo.getAllCalls != 1 {
				t.Fatalf("repo.GetAllUsersRoles called %d times, want 1", repo.getAllCalls)
			}

			if tt.wantErr != nil {
				return
			}

			if !reflect.DeepEqual(got.Roles, tt.repoRoles) {
				t.Errorf("roles = %v, want %v", got.Roles, tt.repoRoles)
			}
		})
	}
}
