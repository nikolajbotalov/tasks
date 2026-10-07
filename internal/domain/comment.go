package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskComment struct {
	ID         uuid.UUID `json:"id"`
	TaskID     uuid.UUID `json:"task_id"`
	AuthorID   uuid.UUID `json:"author_id"`
	AuthorName string    `json:"author_name"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetTaskCommentsResponse struct {
	Comments []TaskComment `json:"comments"`
	Page     int           `json:"page"`
	Limit    int           `json:"limit"`
	Total    int           `json:"total"`
}
