package domain

type TaskPriority string

const (
	TaskLowPriority    TaskPriority = "low"
	TaskMiddlePriority TaskPriority = "middle"
	TaskHighPriority   TaskPriority = "high"
)

func (p TaskPriority) IsValid() bool {
	return p == TaskLowPriority || p == TaskMiddlePriority || p == TaskHighPriority
}
