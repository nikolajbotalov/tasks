package testutils

import "testing"

func AssertOptionalString(t *testing.T, field string, got, want *string) {
	t.Helper()

	switch {
	case got == nil && want == nil:
		return
	case got == nil:
		t.Fatalf("%s = nil, want %q", field, *want)
	case want == nil:
		t.Fatalf("%s = %q, want nil", field, *got)
	case *got != *want:
		t.Errorf("%s = %q, want %q", field, *got, *want)
	}
}
