package domain

type TaskStatus string

const (
	ToDo       TaskStatus = "TO_DO"
	InProgress TaskStatus = "IN_PROGRESS"
	InReview   TaskStatus = "IN_REVIEW"
	Done       TaskStatus = "DONE"
	OnHold     TaskStatus = "ON_HOLD"
	Cancelled  TaskStatus = "CANCELLED"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case ToDo, InProgress, InReview, Done, OnHold, Cancelled:
		return true
	default:
		return false
	}
}
