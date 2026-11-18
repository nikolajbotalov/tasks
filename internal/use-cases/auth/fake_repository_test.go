package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type fakeRepository struct {
	authErr    error
	sessionErr error
	revokeErr  error

	authCalls    int
	sessionCalls int
	revokeCalls  int

	gotAuthID uuid.UUID
	authState domain.AuthState

	gotSessionID uuid.UUID
	sessionState domain.RefreshSessionState

	gotRevokeUser uuid.UUID
}

var _ Repository = (*fakeRepository)(nil)

func (f *fakeRepository) GetAuthState(ctx context.Context, id uuid.UUID) (domain.AuthState, error) {
	f.authCalls++

	f.gotAuthID = id

	return f.authState, f.authErr
}

func (f *fakeRepository) CreateUser(context.Context, string, string, string) error {
	panic("unexpected call: CreateUser")
}

func (f *fakeRepository) FindUserByEmail(context.Context, string) (domain.UserRow, error) {
	panic("unexpected call: FindUserByEmail")
}

func (f *fakeRepository) CreateRefreshSession(context.Context, domain.RefreshSession) error {
	panic("unexpected call: CreateRefreshSession")
}

func (f *fakeRepository) GetRefreshSessionState(ctx context.Context, id uuid.UUID) (domain.RefreshSessionState, error) {
	f.sessionCalls++

	f.gotSessionID = id

	return f.sessionState, f.sessionErr
}

func (f *fakeRepository) RotateRefreshSession(
	context.Context, uuid.UUID, uuid.UUID, uuid.UUID, uuid.UUID, time.Time,
) (bool, domain.RefreshSessionState, error) {
	panic("unexpected call: RotateRefreshSession")
}

func (f *fakeRepository) RevokeRefreshFamilyAndBumpVersion(context.Context, uuid.UUID, uuid.UUID) error {
	panic("unexpected call: RevokeRefreshFamilyAndBumpVersion")
}

func (f *fakeRepository) RevokeAllSessionsAndBumpVersion(ctx context.Context, userID uuid.UUID) error {
	f.revokeCalls++

	f.gotRevokeUser = userID

	return f.revokeErr
}
