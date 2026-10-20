package tasks

import (
	"TaskFlow/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *tasksRepository) CreateTask(id, authorID, assigneeID uuid.UUID, name, description, priority string, status *string) (domain.Task, error) {
	curTime := time.Now()

	var t domain.Task
	err := r.db.QueryRow(`
			INSERT INTO tasks (id, name, description, created_at, updated_at, author_id, assignee_id, status, priority) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
			RETURNING id, name, description, created_at, updated_at, author_id, assignee_id, status, priority`,
		id, name, description, curTime, curTime, authorID, assigneeID, status, priority).
		Scan(&t.ID, &t.Name, &t.Description, &t.CreatedAt, &t.UpdatedAt, &t.AuthorID, &t.AssigneeID, &t.Status, &t.Priority)
	if err != nil {
		return domain.Task{}, fmt.Errorf("insert task %s: %w", id, err)
	}

	return t, nil
}
