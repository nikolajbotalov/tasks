package testutils

import (
	"testing"

	"github.com/google/uuid"
)

func AssertOptionalUUID(t *testing.T, field string, got, want *uuid.UUID) {
	t.Helper()

	switch {
	case got == nil && want == nil:
		return
	case got == nil:
		t.Fatalf("%s = nil, want %s", field, *want)
	case want == nil:
		t.Fatalf("%s = %s, want nil", field, *got)
	case *got != *want:
		t.Errorf("%s = %s, want %s", field, *got, *want)
	}
}
