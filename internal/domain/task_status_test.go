package domain

import "testing"

func TestTaskStatusIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status TaskStatus
		want   bool
	}{
		{name: "to do", status: ToDo, want: true},
		{name: "in progress", status: InProgress, want: true},
		{name: "in review", status: InReview, want: true},
		{name: "done", status: Done, want: true},
		{name: "on hold", status: OnHold, want: true},
		{name: "cancelled", status: Cancelled, want: true},
		{name: "empty string", status: "", want: false},
		{name: "uppercase", status: "TODO", want: false},
		{name: "unknown status", status: "blocked", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.status.IsValid(); got != tt.want {
				t.Errorf("TaskStatus(%q).IsValid() = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}
