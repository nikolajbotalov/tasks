package auth

import (
	"TaskFlow/internal/domain"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("connection refused")

	const password = "password"

	tooLongPassword := strings.Repeat("a", 73)

	tests := []struct {
		name string

		email    string
		password string

		repoErr error

		wantErrs     []error
		wantRepoCall bool
		wantEmail    string
	}{
		{
			name:         "creates user with normalized email and hashed password",
			email:        "  User@Example.COM  ",
			password:     password,
			wantRepoCall: true,
			wantEmail:    "user@example.com",
		},
		{
			name:         "maps email conflict",
			email:        "user@example.com",
			password:     password,
			repoErr:      domain.ErrEmailExists,
			wantErrs:     []error{domain.ErrEmailExists},
			wantRepoCall: true,
			wantEmail:    "user@example.com",
		},
		{
			name:         "wraps unexpected repository error",
			email:        "user@example.com",
			password:     password,
			repoErr:      repoErr,
			wantErrs:     []error{repoErr},
			wantRepoCall: true,
			wantEmail:    "user@example.com",
		},
		{
			name:         "rejects too long password before repository call",
			email:        "user@example.com",
			password:     tooLongPassword,
			wantErrs:     []error{domain.ErrHashedPassword, bcrypt.ErrPasswordTooLong},
			wantRepoCall: false,
			wantEmail:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{createUserErr: tt.repoErr}
			uc := NewAuthUseCases(repo, testConfig())

			gotID, err := uc.CreateUser(t.Context(), tt.email, tt.password)

			if len(tt.wantErrs) == 0 {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			} else {
				for _, want := range tt.wantErrs {
					if !errors.Is(err, want) {
						t.Fatalf("error = %v, want %v", err, want)
					}
				}
			}

			wantCalls := 0
			if tt.wantRepoCall {
				wantCalls = 1
			}

			if repo.createUserCalls != wantCalls {
				t.Fatalf("CreateUser called %d times, want %d", repo.createUserCalls, wantCalls)
			}

			if !tt.wantRepoCall {
				if gotID != uuid.Nil {
					t.Errorf("id = %s, want uuid.Nil", gotID)
				}

				return
			}

			if repo.gotCreateEmail != tt.wantEmail {
				t.Errorf("email = %s, want %s", repo.gotCreateEmail, tt.wantEmail)
			}

			repoID, parseErr := uuid.Parse(repo.gotCreateID)
			if parseErr != nil {
				t.Fatalf("repository id %q is not a uuid: %v", repo.gotCreateID, parseErr)
			}

			if repoID == uuid.Nil {
				t.Errorf("repository id is uuid.Nil")
			}

			if cmpErr := bcrypt.CompareHashAndPassword([]byte(repo.gotCreatePassword), []byte(tt.password)); cmpErr != nil {
				t.Errorf("repository password is not a hash of the input: %v", cmpErr)
			}

			if len(tt.wantErrs) != 0 {
				if gotID != uuid.Nil {
					t.Errorf("id = %s, want uuid.Nil on error", gotID)
				}

				return
			}

			if gotID != repoID {
				t.Errorf("returned id = %s, want %s", gotID, repoID)
			}
		})
	}
}
