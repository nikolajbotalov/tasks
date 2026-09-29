package domain

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type GetAllTasksResponse struct {
	Tasks []Task `json:"tasks"`
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Total int    `json:"total"`
}
