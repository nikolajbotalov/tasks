package tasks

import (
	"TaskFlow/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *tasksRepository) UpdateTask(
	id uuid.UUID,
	assigneeID *uuid.UUID,
	name, description, status, priority *string,
) error {
	res, err := r.db.Exec(
		`UPDATE tasks 
				SET name = COALESCE($1, name), 
				    description = COALESCE($2, description),
				    assignee_id = COALESCE($3, assignee_id),
				    status = COALESCE($4, status),
				    priority = COALESCE($5, priority),
				    updated_at = $6
				WHERE id = $7`, name, description, assigneeID, status, priority, time.Now(), id)
	if err != nil {
		return fmt.Errorf("update task %s: %w", id, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update task %s, affected rows: %w", id, err)
	}

	if affected == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}
