package domain

import "testing"

func TestTaskPriorityIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		priority TaskPriority
		want     bool
	}{
		{name: "low", priority: TaskLowPriority, want: true},
		{name: "middle", priority: TaskMiddlePriority, want: true},
		{name: "high", priority: TaskHighPriority, want: true},
		{name: "empty", priority: "", want: false},
		{name: "upper", priority: "LOW", want: false},
		{name: "unknown", priority: "critical", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.priority.IsValid(); got != tt.want {
				t.Errorf("TaskPriority(%q).IsValid() = %v, want %v", tt.priority, got, tt.want)
			}
		})
	}
}
