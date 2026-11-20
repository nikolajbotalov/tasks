package auth

import (
	"TaskFlow/internal/domain"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestLogout(t *testing.T) {
	t.Parallel()

	var (
		logoutUserID   = uuid.MustParse("11111111-1111-1111-1111-111111111111")
		logoutJTI      = uuid.MustParse("22222222-2222-2222-2222-222222222222")
		logoutFamilyID = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	)

	const ver = 1

	repoErr := errors.New("connection refused")

	validSession := domain.RefreshSessionState{
		UserID:    logoutUserID,
		FamilyID:  logoutFamilyID,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	usedAt := time.Now().Add(-time.Minute)
	revokedAt := time.Now().Add(-time.Minute)

	tests := []struct {
		name string

		token string

		authState domain.AuthState
		authErr   error

		session    domain.RefreshSessionState
		sessionErr error

		revokeErr error
		wantErr   error

		wantAuthCalls    int
		wantSessionCalls int
		wantRevokeCalls  int
	}{
		{
			name:             "revokes sessions on valid token",
			token:            signTestToken(t, logoutUserID, logoutJTI, logoutFamilyID, ver, domain.JWTTypeRefresh, time.Now().Add(time.Hour)),
			authState:        domain.AuthState{Role: domain.RoleUser, TokenVersion: ver},
			session:          validSession,
			wantAuthCalls:    1,
			wantSessionCalls: 1,
			wantRevokeCalls:  1,
		},
		{
			name:    "rejects token with broken signature",
			token:   signTestToken(t, logoutUserID, logoutJTI, logoutFamilyID, ver, domain.JWTTypeRefresh, time.Now().Add(time.Hour)) + "tampered",
			wantErr: domain.ErrInvalidRefreshToken,
		},
		{
			name:    "rejects access token instead of refresh",
			token:   signTestToken(t, logoutUserID, logoutJTI, logoutFamilyID, ver, domain.JWTTypeAccess, time.Now().Add(time.Hour)),
			wantErr: domain.ErrInvalidRefreshToken,
		},
		{
			name:    "rejects expired token",
			token:   signTestToken(t, logoutUserID, logoutJTI, logoutFamilyID, ver, domain.JWTTypeRefresh, time.Now().Add(-time.Minute)),
			wantErr: domain.ErrInvalidRefreshToken,
		},
		{
			name:      "rejects token with stale version",
			token:     signTestToken(t, logoutUserID, logoutJTI, logoutFamilyID, ver+1, domain.JWTTypeRefresh, time.Now().Add(time.Hour)),
			authState: domain.AuthState{Role: domain.RoleUser, TokenVersion: ver},
			wantErr:   domain.ErrInvalidRefreshToken,

			wantAuthCalls: 1,
		},
		{
			name:      "returns user not found",
			token:     signTestToken(t, logoutUserID, logoutJTI, logoutFamilyID, ver, domain.JWTTypeRefresh, time.Now().Add(time.Hour)),
			authState: domain.AuthState{Role: domain.RoleUser, TokenVersion: ver},
			authErr:   domain.ErrUserNotFound,
			wantErr:   domain.ErrUserNotFound,

			wantAuthCalls: 1,
		},
		{
			name:    "wraps unexpected auth state error",
			token:   signTestToken(t, logoutUserID, logoutJTI, logoutFamilyID, ver, domain.JWTTypeRefresh, time.Now().Add(time.Hour)),
			authErr: repoErr,
			wantErr: repoErr,

			wantAuthCalls: 1,
		},
		{
			name:       "returns invalid token when session is missing",
			token:      signTestToken(t, logoutUserID, logoutJTI, logoutFamilyID, ver, domain.JWTTypeRefresh, time.Now().Add(time.Hour)),
			authState:  domain.AuthState{Role: domain.RoleUser, TokenVersion: ver},
			sessionErr: domain.ErrInvalidRefreshToken,
			wantErr:    domain.ErrInvalidRefreshToken,

			wantAuthCalls:    1,
			wantSessionCalls: 1,
		},
		{
			name:      "rejects session of another user",
			token:     signTestToken(t, logoutUserID, logoutJTI, logoutFamilyID, ver, domain.JWTTypeRefresh, time.Now().Add(time.Hour)),
			authState: domain.AuthState{Role: domain.RoleUser, TokenVersion: ver},
			session:   domain.RefreshSessionState{UserID: uuid.New(), FamilyID: logoutFamilyID, ExpiresAt: time.Now().Add(time.Hour)},
			wantErr:   domain.ErrInvalidRefreshToken,

			wantAuthCalls:    1,
			wantSessionCalls: 1,
		},
		{
			name:      "rejects already used session",
			token:     signTestToken(t, logoutUserID, logoutJTI, logoutFamilyID, ver, domain.JWTTypeRefresh, time.Now().Add(time.Hour)),
			authState: domain.AuthState{Role: domain.RoleUser, TokenVersion: ver},
			session:   domain.RefreshSessionState{UserID: logoutUserID, FamilyID: logoutFamilyID, UsedAt: &usedAt, ExpiresAt: time.Now().Add(time.Hour)},
			wantErr:   domain.ErrInvalidRefreshToken,

			wantAuthCalls:    1,
			wantSessionCalls: 1,
		},
		{
			name:      "rejects revoked session",
			token:     signTestToken(t, logoutUserID, logoutJTI, logoutFamilyID, ver, domain.JWTTypeRefresh, time.Now().Add(time.Hour)),
			authState: domain.AuthState{Role: domain.RoleUser, TokenVersion: ver},
			session:   domain.RefreshSessionState{UserID: logoutUserID, FamilyID: logoutFamilyID, RevokedAt: &revokedAt, ExpiresAt: time.Now().Add(time.Hour)},
			wantErr:   domain.ErrInvalidRefreshToken,

			wantAuthCalls:    1,
			wantSessionCalls: 1,
		},
		{
			name:      "rejects expired session",
			token:     signTestToken(t, logoutUserID, logoutJTI, logoutFamilyID, ver, domain.JWTTypeRefresh, time.Now().Add(time.Hour)),
			authState: domain.AuthState{Role: domain.RoleUser, TokenVersion: ver},
			session:   domain.RefreshSessionState{UserID: logoutUserID, FamilyID: logoutFamilyID, ExpiresAt: time.Now().Add(-time.Minute)},
			wantErr:   domain.ErrInvalidRefreshToken,

			wantAuthCalls:    1,
			wantSessionCalls: 1,
		},
		{
			name:      "wraps revoke error",
			token:     signTestToken(t, logoutUserID, logoutJTI, logoutFamilyID, ver, domain.JWTTypeRefresh, time.Now().Add(time.Hour)),
			authState: domain.AuthState{Role: domain.RoleUser, TokenVersion: ver},
			session:   validSession,
			revokeErr: repoErr,
			wantErr:   repoErr,

			wantAuthCalls:    1,
			wantSessionCalls: 1,
			wantRevokeCalls:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{
				authState:    tt.authState,
				authErr:      tt.authErr,
				sessionState: tt.session,
				sessionErr:   tt.sessionErr,
				revokeErr:    tt.revokeErr,
			}
			uc := NewAuthUseCases(repo, testConfig())

			err := uc.Logout(t.Context(), tt.token)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if repo.authCalls != tt.wantAuthCalls {
				t.Errorf("GetAuthState called %d times, want %d", repo.authCalls, tt.wantAuthCalls)
			}

			if repo.sessionCalls != tt.wantSessionCalls {
				t.Errorf("GetRefreshSessionState called %d times, want %d", repo.sessionCalls, tt.wantSessionCalls)
			}

			if repo.revokeCalls != tt.wantRevokeCalls {
				t.Errorf("RevokeAllSessionsAndBumpVersion called %d times, want %d", repo.revokeCalls, tt.wantRevokeCalls)
			}

			if tt.wantRevokeCalls == 1 && repo.gotRevokeUser != logoutUserID {
				t.Errorf("revoked userID = %s, want %s", repo.gotRevokeUser, logoutUserID)
			}
		})
	}
}
