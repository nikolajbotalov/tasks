package tasks

import (
	"TaskFlow/internal/domain"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *tasksRepository) GetTaskByID(id uuid.UUID) (domain.Task, error) {
	var t domain.Task

	err := r.db.QueryRow(`
		SELECT id, name, description, created_at, updated_at, author_id, assignee_id, status, priority 
		FROM tasks WHERE id = $1`, id).
		Scan(&t.ID, &t.Name, &t.Description, &t.CreatedAt, &t.UpdatedAt, &t.AuthorID, &t.AssigneeID, &t.Status, &t.Priority)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Task{}, domain.ErrTaskNotFound
		}

		return domain.Task{}, fmt.Errorf("get task %s: %w", id.String(), err)
	}

	t.Status = TaskStatusToLower(t.Status)
	return t, nil
}
