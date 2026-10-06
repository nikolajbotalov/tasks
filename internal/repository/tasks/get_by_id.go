package tasks

import (
	"TaskFlow/internal/domain"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func (r *tasksRepository) GetTaskByID(id uuid.UUID) (domain.Task, error) {
	var task domain.Task

	err := r.db.QueryRow("SELECT id, name, description, created_at, updated_at, author_id FROM tasks WHERE id = $1", id).
		Scan(&task.ID, &task.Name, &task.Description, &task.CreatedAt, &task.UpdatedAt, &task.AuthorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Task{}, domain.ErrTaskNotFound
		}

		return domain.Task{}, fmt.Errorf("get task %s: %w", id.String(), err)
	}

	return task, nil
}
