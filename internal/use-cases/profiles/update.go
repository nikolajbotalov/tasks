package profiles

import (
	"TaskFlow/internal/domain"
	"context"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	profileFieldMinLength = 2
	profileFieldMaxLength = 100
)

func (uc *profilesUseCases) UpdateProfile(
	ctx context.Context,
	username, firstName, lastName *string,
	userID, callerID uuid.UUID,
) error {
	if callerID != userID {
		return domain.ErrForeignProfileChange
	}

	normalizedUserName, err := normalizeOptionalField(username)
	if err != nil {
		return err
	}

	normalizedFirstName, err := normalizeOptionalField(firstName)
	if err != nil {
		return err
	}

	normalizedLastName, err := normalizeOptionalField(lastName)
	if err != nil {
		return err
	}

	return uc.repo.UpdateProfile(ctx, normalizedUserName, normalizedFirstName, normalizedLastName, userID)
}

func normalizeOptionalField(val *string) (*string, error) {
	if val == nil {
		return nil, nil
	}

	trimmed := strings.TrimSpace(*val)
	length := utf8.RuneCountInString(trimmed)
	if length < profileFieldMinLength || length > profileFieldMaxLength {
		return nil, domain.ErrInvalidProfileData
	}

	return &trimmed, nil
}
