package profiles

import (
	"TaskFlow/internal/domain"
	"errors"
	"strings"
	"testing"
)

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
