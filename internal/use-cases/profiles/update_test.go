package profiles

import (
	"TaskFlow/internal/domain"
	"TaskFlow/internal/testutils"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
)

var (
	profileUserID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	otherUserID   = uuid.MustParse("22222222-2222-2222-2222-222222222222")
)

func TestUpdateProfile(t *testing.T) {
	t.Parallel()

	ptr := func(s string) *string { return &s }

	tests := []struct {
		name string

		username  *string
		firstName *string
		lastName  *string
		callerID  uuid.UUID

		repoErr error
		wantErr error

		wantUserName  *string
		wantFirstName *string
		wantLastName  *string
		wantRepoCall  bool
	}{
		{
			name:     "foreign profile",
			username: ptr("user1"),
			callerID: otherUserID,
			wantErr:  domain.ErrForeignProfileChange,
		},
		{
			name:     "foreign profile wins over invalid field",
			username: ptr("user1"),
			callerID: otherUserID,
			wantErr:  domain.ErrForeignProfileChange,
		},
		{
			name:         "all fields nil",
			callerID:     profileUserID,
			wantRepoCall: true,
		},
		{
			name:          "trim fields",
			username:      ptr("  user1  "),
			firstName:     ptr("  test  "),
			lastName:      ptr("  test  "),
			callerID:      profileUserID,
			wantRepoCall:  true,
			wantUserName:  ptr("user1"),
			wantFirstName: ptr("test"),
			wantLastName:  ptr("test"),
		},
		{
			name:     "invalid username",
			username: ptr("u"),
			callerID: profileUserID,
			wantErr:  domain.ErrInvalidProfileData,
		},
		{
			name:         "repo error is propagated",
			username:     ptr("user1"),
			callerID:     profileUserID,
			repoErr:      domain.ErrProfileNotFound,
			wantErr:      domain.ErrProfileNotFound,
			wantRepoCall: true,
			wantUserName: ptr("user1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{updateErr: tt.repoErr}
			uc := NewProfilesUseCases(repo)

			err := uc.UpdateProfile(t.Context(), tt.username, tt.firstName, tt.lastName, profileUserID, tt.callerID)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if !tt.wantRepoCall {
				if repo.updateCalls != 0 {
					t.Fatalf("repo.UpdateProfile called %d times, want 0", repo.updateCalls)
				}
				return
			}

			if repo.updateCalls != 1 {
				t.Fatalf("repo.UpdateProfile called %d times, want 1", repo.updateCalls)
			}

			if repo.gotUserID != profileUserID {
				t.Errorf("userID = %s, want %s", repo.gotUserID, profileUserID)
			}

			testutils.AssertOptionalString(t, "username", repo.gotUsername, tt.wantUserName)
			testutils.AssertOptionalString(t, "firstName", repo.gotFirstName, tt.wantFirstName)
			testutils.AssertOptionalString(t, "lastName", repo.gotLastName, tt.wantLastName)
		})
	}
}

func TestNormalizeOptionalField(t *testing.T) {
	t.Parallel()

	ptr := func(s string) *string { return &s }

	tests := []struct {
		name    string
		input   *string
		want    *string
		wantErr error
	}{
		{name: "nil stays nil", input: nil, want: nil},
		{name: "single rune", input: ptr("A"), wantErr: domain.ErrInvalidProfileData},
		{name: "empty string", input: ptr(""), wantErr: domain.ErrInvalidProfileData},
		{name: "only spaces", input: ptr("  "), wantErr: domain.ErrInvalidProfileData},
		{name: "min length two runes", input: ptr("ab"), want: ptr("ab")},
		{name: "trims surrounding spaces", input: ptr("  Test  "), want: ptr("Test")},
		{name: "cyrillic two runes", input: ptr("Те"), want: ptr("Те")},
		{name: "max length 100 runes", input: ptr(strings.Repeat("a", 100)), want: ptr(strings.Repeat("a", 100))},
		{name: "101 runes", input: ptr(strings.Repeat("a", 101)), wantErr: domain.ErrInvalidProfileData},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeOptionalField(tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if tt.want == nil {
				if got != nil {
					t.Fatalf("got %q, want nil", *got)
				}
				return
			}

			if got == nil {
				t.Fatalf("got nil, want %q", *tt.want)
			}

			if *got != *tt.want {
				t.Errorf("get %q, want %q", *got, *tt.want)
			}
		})
	}
}
