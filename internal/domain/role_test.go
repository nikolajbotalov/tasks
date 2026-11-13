package domain

import "testing"

func TestRoleIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		role Role
		want bool
	}{
		{name: "viewer", role: RoleViewer, want: true},
		{name: "user", role: RoleUser, want: true},
		{name: "admin", role: RoleAdmin, want: true},
		{name: "empty string", role: "", want: false},
		{name: "unknown role", role: "owner", want: false},
		{name: "uppercase", role: "Admin", want: false},
		{name: "leading spaces", role: " admin", want: false},
		{name: "trailing spaces", role: "admin ", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.role.IsValid(); got != tt.want {
				t.Errorf("Role(%q).IsValid() = %v, want %v", tt.role, got, tt.want)
			}
		})
	}
}
