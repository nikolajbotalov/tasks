package auth

import (
	"TaskFlow/internal/domain"
	"context"
	"time"

	"github.com/google/uuid"
)

type fakeRepository struct {
	authErr          error
	sessionErr       error
	revokeErr        error
	rotateErr        error
	revokeFamilyErr  error
	findErr          error
	createSessionErr error

	authCalls          int
	sessionCalls       int
	revokeCalls        int
	rotateCalls        int
	revokeFamilyCalls  int
	findCalls          int
	createSessionCalls int

	gotAuthID uuid.UUID
	authState domain.AuthState

	gotSessionID uuid.UUID
	sessionState domain.RefreshSessionState

	gotRevokeUser uuid.UUID

	rotated    bool
	oldSession domain.RefreshSessionState

	gotRotateOldID   uuid.UUID
	gotRotateNewID   uuid.UUID
	gotRotateUserID  uuid.UUID
	gotRotateFamily  uuid.UUID
	gotRotateExpires time.Time

	gotRevokeFamilyUser uuid.UUID
	gotRevokeFamilyID   uuid.UUID

	userRow      domain.UserRow
	gotFindEmail string

	gotSession domain.RefreshSession
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

func (f *fakeRepository) FindUserByEmail(ctx context.Context, email string) (domain.UserRow, error) {
	f.findCalls++

	f.gotFindEmail = email

	return f.userRow, f.findErr
}

func (f *fakeRepository) CreateRefreshSession(ctx context.Context, session domain.RefreshSession) error {
	f.createSessionCalls++

	f.gotSession = session

	return f.createSessionErr
}

func (f *fakeRepository) GetRefreshSessionState(ctx context.Context, id uuid.UUID) (domain.RefreshSessionState, error) {
	f.sessionCalls++

	f.gotSessionID = id

	return f.sessionState, f.sessionErr
}

func (f *fakeRepository) RotateRefreshSession(
	ctx context.Context,
	oldID, newID, userID, familyID uuid.UUID,
	expiresAt time.Time,
) (bool, domain.RefreshSessionState, error) {
	f.rotateCalls++

	f.gotRotateOldID = oldID
	f.gotRotateNewID = newID
	f.gotRotateUserID = userID
	f.gotRotateFamily = familyID
	f.gotRotateExpires = expiresAt

	return f.rotated, f.oldSession, f.rotateErr
}

func (f *fakeRepository) RevokeRefreshFamilyAndBumpVersion(
	ctx context.Context,
	userID, familyID uuid.UUID,
) error {
	f.revokeFamilyCalls++

	f.gotRevokeFamilyUser = userID
	f.gotRevokeFamilyID = familyID

	return f.revokeFamilyErr
}

func (f *fakeRepository) RevokeAllSessionsAndBumpVersion(ctx context.Context, userID uuid.UUID) error {
	f.revokeCalls++

	f.gotRevokeUser = userID

	return f.revokeErr
}
