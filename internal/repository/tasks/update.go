package tasks

import (
	"TaskFlow/internal/domain"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *tasksRepository) UpdateTask(id uuid.UUID, name, description *string) error {
	res, err := r.db.Exec(
		`UPDATE tasks 
				SET name = COALESCE($1, name), 
				    description = COALESCE($2, description), 
				    updated_at = $3
				WHERE id = $4`, name, description, time.Now(), id)
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
