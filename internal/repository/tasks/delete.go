package tasks

import (
	"TaskFlow/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

func (r *tasksRepository) DeleteTask(id, authorID uuid.UUID) error {
	res, err := r.db.Exec("DELETE FROM tasks WHERE id = $1 AND author_id = $2", id, authorID)
	if err != nil {
		return fmt.Errorf("delete task %s: %w", id, err)
	}

	var affected int64
	affected, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete task %s, rows affected: %w", id, err)
	}

	if affected == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}
