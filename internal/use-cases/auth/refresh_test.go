package auth

import (
	"TaskFlow/internal/domain"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRefreshToken(t *testing.T) {
	t.Parallel()

	var (
		refreshUserID   = uuid.MustParse("11111111-1111-1111-1111-111111111111")
		refreshJTI      = uuid.MustParse("22222222-2222-2222-2222-222222222222")
		refreshFamilyID = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	)

	const ver = 1

	repoErr := errors.New("connection refused")

	validAuthState := domain.AuthState{Role: domain.RoleUser, TokenVersion: ver}

	validToken := signTestToken(
		t,
		refreshUserID, refreshJTI, refreshFamilyID,
		ver,
		domain.JWTTypeRefresh,
		time.Now().Add(time.Hour),
	)

	usedWithinGrace := time.Now().Add(-5 * time.Second)
	usedLongAgo := time.Now().Add(-time.Minute)
	revokedAt := time.Now().Add(-time.Minute)

	matchingSession := domain.RefreshSessionState{
		UserID:    refreshUserID,
		FamilyID:  refreshFamilyID,
		UsedAt:    &usedLongAgo,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	tests := []struct {
		name string

		token string

		authState domain.AuthState
		authErr   error

		rotated    bool
		oldSession domain.RefreshSessionState
		rotateErr  error

		revokeFamilyErr error

		wantErr error

		wantAuthCalls         int
		wantRotateCalls       int
		wantRevokeFamilyCalls int
	}{
		{
			name:      "rotates token pair",
			token:     validToken,
			authState: validAuthState,
			rotated:   true,

			wantAuthCalls:   1,
			wantRotateCalls: 1,
		},
		{
			name:    "rejects token with broken signature",
			token:   validToken + "tampered",
			wantErr: domain.ErrInvalidRefreshToken,
		},
		{
			name:    "rejects access token instead of refresh",
			token:   signTestToken(t, refreshUserID, refreshJTI, refreshFamilyID, ver, domain.JWTTypeAccess, time.Now().Add(time.Hour)),
			wantErr: domain.ErrInvalidRefreshToken,
		},
		{
			name:    "rejects expired token",
			token:   signTestToken(t, refreshUserID, refreshJTI, refreshFamilyID, ver, domain.JWTTypeRefresh, time.Now().Add(-time.Minute)),
			wantErr: domain.ErrInvalidRefreshToken,
		},
		{
			name:    "maps user not found to invalid token",
			token:   validToken,
			authErr: domain.ErrUserNotFound,
			wantErr: domain.ErrInvalidRefreshToken,

			wantAuthCalls: 1,
		},
		{
			name:    "wraps unexpected auth state error",
			token:   validToken,
			authErr: repoErr,
			wantErr: repoErr,

			wantAuthCalls: 1,
		},
		{
			name:      "rejects token with stale version",
			token:     signTestToken(t, refreshUserID, refreshJTI, refreshFamilyID, ver+1, domain.JWTTypeRefresh, time.Now().Add(time.Hour)),
			authState: validAuthState,
			wantErr:   domain.ErrInvalidRefreshToken,

			wantAuthCalls: 1,
		},
		{
			name:      "wraps rotation error",
			token:     validToken,
			authState: validAuthState,
			rotateErr: repoErr,
			wantErr:   repoErr,

			wantAuthCalls:   1,
			wantRotateCalls: 1,
		},
		{
			name:      "rejects session of another user",
			token:     validToken,
			authState: validAuthState,
			oldSession: domain.RefreshSessionState{
				UserID:    uuid.New(),
				FamilyID:  refreshFamilyID,
				UsedAt:    &usedLongAgo,
				ExpiresAt: time.Now().Add(time.Hour),
			},
			wantErr: domain.ErrInvalidRefreshToken,

			wantAuthCalls:   1,
			wantRotateCalls: 1,
		},
		{
			name:      "rejects session from another family",
			token:     validToken,
			authState: validAuthState,
			oldSession: domain.RefreshSessionState{
				UserID:    refreshUserID,
				FamilyID:  uuid.New(),
				UsedAt:    &usedLongAgo,
				ExpiresAt: time.Now().Add(time.Hour),
			},
			wantErr: domain.ErrInvalidRefreshToken,

			wantAuthCalls:   1,
			wantRotateCalls: 1,
		},
		{
			name:      "rejects revoked session",
			token:     validToken,
			authState: validAuthState,
			oldSession: domain.RefreshSessionState{
				UserID:    refreshUserID,
				FamilyID:  refreshFamilyID,
				RevokedAt: &revokedAt,
				UsedAt:    &usedLongAgo,
				ExpiresAt: time.Now().Add(time.Hour),
			},
			wantErr: domain.ErrInvalidRefreshToken,

			wantAuthCalls:   1,
			wantRotateCalls: 1,
		},
		{
			name:      "rejects expired session",
			token:     validToken,
			authState: validAuthState,
			oldSession: domain.RefreshSessionState{
				UserID:    refreshUserID,
				FamilyID:  refreshFamilyID,
				ExpiresAt: time.Now().Add(-time.Minute),
			},
			wantErr: domain.ErrInvalidRefreshToken,

			wantAuthCalls:   1,
			wantRotateCalls: 1,
		},
		{
			name:      "rejects session that was never used",
			token:     validToken,
			authState: validAuthState,
			oldSession: domain.RefreshSessionState{
				UserID:    refreshUserID,
				FamilyID:  refreshFamilyID,
				ExpiresAt: time.Now().Add(time.Hour),
			},
			wantErr: domain.ErrInvalidRefreshToken,

			wantAuthCalls:   1,
			wantRotateCalls: 1,
		},
		{
			name:      "returns refresh in progress inside grace period",
			token:     validToken,
			authState: validAuthState,
			oldSession: domain.RefreshSessionState{
				UserID:    refreshUserID,
				FamilyID:  refreshFamilyID,
				UsedAt:    &usedWithinGrace,
				ExpiresAt: time.Now().Add(time.Hour),
			},
			wantErr: domain.ErrRefreshInProgress,

			wantAuthCalls:   1,
			wantRotateCalls: 1,
		},
		{
			name:       "revokes family on reuse after grace period",
			token:      validToken,
			authState:  validAuthState,
			oldSession: matchingSession,
			wantErr:    domain.ErrInvalidRefreshToken,

			wantAuthCalls:         1,
			wantRotateCalls:       1,
			wantRevokeFamilyCalls: 1,
		},
		{
			name:       "maps user not found on family revoke",
			token:      validToken,
			authState:  validAuthState,
			oldSession: matchingSession,

			revokeFamilyErr: domain.ErrUserNotFound,
			wantErr:         domain.ErrInvalidRefreshToken,

			wantAuthCalls:         1,
			wantRotateCalls:       1,
			wantRevokeFamilyCalls: 1,
		},
		{
			name:       "wraps family revoke error",
			token:      validToken,
			authState:  validAuthState,
			oldSession: matchingSession,

			revokeFamilyErr: repoErr,
			wantErr:         repoErr,

			wantAuthCalls:         1,
			wantRotateCalls:       1,
			wantRevokeFamilyCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := &fakeRepository{
				authState:       tt.authState,
				authErr:         tt.authErr,
				rotated:         tt.rotated,
				oldSession:      tt.oldSession,
				rotateErr:       tt.rotateErr,
				revokeFamilyErr: tt.revokeFamilyErr,
			}
			uc := NewAuthUseCases(repo, testConfig())

			got, err := uc.RefreshToken(t.Context(), tt.token)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}

			if repo.authCalls != tt.wantAuthCalls {
				t.Errorf("GetAuthState called %d times, want %d", repo.authCalls, tt.wantAuthCalls)
			}

			if repo.rotateCalls != tt.wantRotateCalls {
				t.Errorf("RotateRefreshSession called %d times, want %d", repo.rotateCalls, tt.wantRotateCalls)
			}

			if repo.revokeFamilyCalls != tt.wantRevokeFamilyCalls {
				t.Errorf("RevokeRefreshFamilyAndBumpVersion called %d times, want %d", repo.revokeFamilyCalls, tt.wantRevokeFamilyCalls)
			}

			if tt.wantRotateCalls == 1 {
				if repo.gotRotateOldID != refreshJTI {
					t.Errorf("rotated odlID = %s, want %s", repo.gotRotateOldID, refreshJTI)
				}

				if repo.gotRotateNewID == uuid.Nil || repo.gotRotateNewID == refreshJTI {
					t.Errorf("rotated newID = %s, want new non-nil id", repo.gotRotateNewID)
				}

				if repo.gotRotateUserID != refreshUserID {
					t.Errorf("rotated userID = %s, want %s", repo.gotRotateUserID, refreshUserID)
				}

				if repo.gotRotateFamily != refreshFamilyID {
					t.Errorf("rotated family = %s, want %s", repo.gotRotateFamily, refreshFamilyID)
				}

				if repo.gotRotateExpires.IsZero() {
					t.Errorf("rotated expiresAt is zero")
				}
			}

			if tt.wantRevokeFamilyCalls == 1 {
				if repo.gotRevokeFamilyUser != refreshUserID {
					t.Errorf("revoked family userID = %s, want %s", repo.gotRevokeFamilyUser, refreshUserID)
				}

				if repo.gotRevokeFamilyID != refreshFamilyID {
					t.Errorf("revoked familyID = %s, want %s", repo.gotRevokeFamilyID, refreshFamilyID)
				}
			}

			if tt.wantErr != nil {
				if got != (domain.TokenPair{}) {
					t.Errorf("token paid = %+v, want zero pai on error", got)
				}

				return
			}

			if got.AccessToken == "" || got.RefreshToken == "" {
				t.Errorf("token pair is incomplete: %+v", got)
			}

			if !got.RefreshExpiresAt.After(time.Now()) {
				t.Errorf("RefreshExpiresAt = %s, want in the future", got.RefreshExpiresAt)
			}
		})
	}
}
