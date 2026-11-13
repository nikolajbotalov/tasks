package roles

import (
	"TaskFlow/internal/domain"
	"errors"
	"testing"

	"github.com/google/uuid"
)

var (
	testTargetID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	testCallerID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
)

func TestUpdateUserRole(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("connection refused")

	tests := []struct {
		name string

		targetID uuid.UUID
		callerID uuid.UUID
		role     string
		repoErr  error

		wantErr      error
		wantRepoCall bool
		wantRole     domain.Role
	}{
		{
			name:         "assigns valid role",
			targetID:     testTargetID,
			callerID:     testCallerID,
			role:         string(domain.RoleAdmin),
			wantRepoCall: true,
			wantRole:     domain.RoleAdmin,
		},
		{
			name:     "unknown role is rejected",
			targetID: testTargetID,
			callerID: testCallerID,
			role:     "owner",
			wantErr:  domain.ErrUnknownRole,
		},
		{
			name:     "role is case sensitive",
			targetID: testTargetID,
			callerID: testCallerID,
			role:     "Admin",
			wantErr:  domain.ErrUnknownRole,
		},
		{
			name:     "own role change is forbidden",
			targetID: testCallerID,
			callerID: testCallerID,
			role:     string(domain.RoleViewer),
			wantErr:  domain.ErrSelfRoleChange,
		},
		{
			name:     "unknown role wins over own role change",
			targetID: testCallerID,
			callerID: testCallerID,
			role:     "owner",
			wantErr:  domain.ErrUnknownRole,
		},
		{
			name:         "repo error is propagated",
			targetID:     testTargetID,
			callerID:     testCallerID,
			role:         string(domain.RoleUser),
			repoErr:      repoErr,
			wantErr:      repoErr,
			wantRepoCall: true,
			wantRole:     domain.RoleUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{updateErr: tt.repoErr}
			uc := NewRolesUseCases(repo)

			err := uc.UpdateUserRole(t.Context(), tt.targetID, tt.callerID, tt.role)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantRole)
			}

			if !tt.wantRepoCall {
				if repo.updateCalls != 0 {
					t.Fatalf("repo.UpdateUserRole called %d times, want 0", repo.updateCalls)
				}
				return
			}

			if repo.updateCalls != 1 {
				t.Fatalf("repo.UpdateUserRole called %d times, want 1", repo.updateCalls)
			}

			if repo.gotTargetID != tt.targetID {
				t.Errorf("targetID = %s, want %s", repo.gotTargetID, tt.targetID)
			}

			if repo.gotRole != tt.wantRole {
				t.Errorf("role = %s, want %s", repo.gotRole, tt.wantRole)
			}
		})
	}
}
