package domain

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AuthorID    uuid.UUID `json:"author_id"`
	AssigneeID  uuid.UUID `json:"assignee_id"`
	Status      *string   `json:"status"`
	Priority    *string   `json:"priority"`
}

type GetAllTasksResponse struct {
	Tasks []Task `json:"tasks"`
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Total int    `json:"total"`
}
