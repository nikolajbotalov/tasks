package tasks

import (
	"TaskFlow/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *tasksRepository) CreateTask(id, authorID uuid.UUID, name, description string) (domain.Task, error) {
	curTime := time.Now()

	var task domain.Task
	err := r.db.QueryRow("INSERT INTO tasks (id, name, description, created_at, updated_at, author_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, description, created_at, updated_at, author_id",
		id, name, description, curTime, curTime, authorID).Scan(&task.ID, &task.Name, &task.Description, &task.CreatedAt, &task.UpdatedAt, &task.AuthorID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("insert task %s: %w", id, err)
	}

	return task, nil
}
