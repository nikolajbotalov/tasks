package auth

import (
	"TaskFlow/internal/domain"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUser(t *testing.T) {
	t.Parallel()

	var (
		loginUserID   = uuid.MustParse("11111111-1111-1111-1111-111111111111")
		loginEmail    = "user@example.com"
		loginPassword = "secret"
	)

	const ver = 1

	repoErr := errors.New("connection refused")

	validHash, err := bcrypt.GenerateFromPassword([]byte(loginPassword), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash test password: %v", err)
	}

	validUser := domain.UserRow{
		ID:           loginUserID,
		Password:     string(validHash),
		TokenVersion: ver,
	}

	tests := []struct {
		name string

		email    string
		password string

		userRow domain.UserRow
		findErr error

		createSessionErr error

		wantErr       error
		wantFindEmail string
		wantUserID    uuid.UUID
		wantTokens    bool
		wantSession   bool
	}{
		{
			name:          "logs in with valid credentials",
			email:         "  User@Example.COM  ",
			password:      loginPassword,
			userRow:       validUser,
			wantFindEmail: loginEmail,
			wantUserID:    loginUserID,
			wantTokens:    true,
			wantSession:   true,
		},
		{
			name:          "maps user not found to invalid credentials",
			email:         loginEmail,
			password:      loginPassword,
			findErr:       domain.ErrUserNotFound,
			wantErr:       domain.ErrInvalidCredentials,
			wantFindEmail: loginEmail,
		},
		{
			name:          "maps wrong password to invalid credentials",
			email:         loginEmail,
			password:      "wrong-password",
			userRow:       validUser,
			wantErr:       domain.ErrInvalidCredentials,
			wantFindEmail: loginEmail,
		},
		{
			name:          "wraps unexpected repository error",
			email:         loginEmail,
			password:      loginPassword,
			findErr:       repoErr,
			wantErr:       repoErr,
			wantFindEmail: loginEmail,
		},
		{
			name:          "wraps unexpected password compare error",
			email:         loginEmail,
			password:      loginPassword,
			userRow:       domain.UserRow{ID: loginUserID, Password: "not-a-bcrypt-hash", TokenVersion: ver},
			wantErr:       bcrypt.ErrHashTooShort,
			wantFindEmail: loginEmail,
		},
		{
			name:             "wraps refresh session creation error",
			email:            loginEmail,
			password:         loginPassword,
			userRow:          validUser,
			createSessionErr: repoErr,
			wantErr:          repoErr,
			wantFindEmail:    loginEmail,
			wantSession:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{
				userRow:          tt.userRow,
				findErr:          tt.findErr,
				createSessionErr: tt.createSessionErr,
			}
			uc := NewAuthUseCases(repo, testConfig())

			gotTokens, gotUserID, err := uc.LoginUser(t.Context(), tt.email, tt.password)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if repo.findCalls != 1 {
				t.Fatalf("FindUserByEmail called %d times, want 1", repo.findCalls)
			}

			if repo.gotFindEmail != tt.wantFindEmail {
				t.Errorf("email = %q, want = %q", repo.gotFindEmail, tt.wantFindEmail)
			}

			wantSessionCalls := 0
			if tt.wantSession {
				wantSessionCalls = 1
			}

			if repo.createSessionCalls != wantSessionCalls {
				t.Errorf("CreateRefreshSession called %d times, want %d", repo.createSessionErr, wantSessionCalls)
			}

			if tt.wantSession {
				if repo.gotSession.UserID != tt.userRow.ID {
					t.Errorf("session.ID = %s, want %s", repo.gotSession.UserID, tt.userRow.ID)
				}

				if repo.gotSession.ID == uuid.Nil {
					t.Errorf("session.ID is nil")
				}

				if repo.gotSession.FamilyID == uuid.Nil {
					t.Errorf("session.FamilyID is nil")
				}

				if !repo.gotSession.ExpiresAt.After(time.Now()) {
					t.Errorf("session.ExpiresAt = %s, want in the future", repo.gotSession.ExpiresAt)
				}
			}

			if tt.wantErr != nil {
				if gotTokens != (domain.TokenPair{}) {
					t.Errorf("token pair = %+v, want zero pair on error", gotTokens)
				}

				if gotUserID != uuid.Nil {
					t.Errorf("userID = %s, want uuid.Nil on error", gotUserID)
				}

				return
			}

			if !tt.wantTokens {
				t.Fatalf("test case expects tokens but wantTokens is false")
			}

			if gotTokens.AccessToken == "" || gotTokens.RefreshToken == "" {
				t.Errorf("token pair is incomplete: %+v", gotTokens)
			}

			if !gotTokens.RefreshExpiresAt.After(time.Now()) {
				t.Errorf("RefreshExpiresAt = %s, want in the future", gotTokens.RefreshExpiresAt)
			}

			if gotUserID != tt.wantUserID {
				t.Errorf("userID = %s, want %s", gotUserID, tt.wantUserID)
			}
		})
	}
}
