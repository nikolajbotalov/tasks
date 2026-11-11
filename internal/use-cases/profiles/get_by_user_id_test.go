package profiles

import (
	"TaskFlow/internal/domain"
	"TaskFlow/internal/testutils"
	"errors"
	"testing"
)

func TestGetProfileByUserID(t *testing.T) {
	t.Parallel()

	ptr := func(s string) *string { return &s }

	requestedID := otherUserID

	repoProfile := domain.Profile{
		UserID:   requestedID,
		Username: ptr("user"),
	}

	tests := []struct {
		name        string
		repoErr     error
		repoProfile domain.Profile
		wantErr     error
	}{
		{
			name:        "return profile from repository",
			repoProfile: repoProfile,
		},
		{
			name:    "user not found is propagated",
			repoErr: domain.ErrUserNotFound,
			wantErr: domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{
				getErr:  tt.repoErr,
				profile: tt.repoProfile,
			}
			uc := NewProfilesUseCases(repo)

			got, err := uc.GetProfileByUserID(t.Context(), requestedID)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if repo.getCalls != 1 {
				t.Fatalf("repo.GetProfileByUserID called %d times, want 1", repo.getCalls)
			}

			if repo.gotGetUserID != requestedID {
				t.Errorf("userID = %s, want %s", got.UserID, requestedID)
			}

			if got.UserID != tt.repoProfile.UserID {
				t.Errorf("profile.UserID = %s, want %s", got.UserID, tt.repoProfile.UserID)
			}

			testutils.AssertOptionalString(t, "username", got.Username, tt.repoProfile.Username)
		})
	}
}
