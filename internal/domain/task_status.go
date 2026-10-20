package domain

type TaskStatus string

const (
	ToDo       TaskStatus = "to_do"
	InProgress TaskStatus = "in_progress"
	InReview   TaskStatus = "in_review"
	Done       TaskStatus = "done"
	OnHold     TaskStatus = "on_hold"
	Cancelled  TaskStatus = "cancelled"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case ToDo, InProgress, InReview, Done, OnHold, Cancelled:
		return true
	default:
		return false
	}
}
