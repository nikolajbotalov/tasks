package tasks

import (
	"TaskFlow/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

func (r *tasksRepository) GetTaskList(authorID uuid.UUID, page, limit int) ([]domain.Task, int, error) {
	rows, err := r.db.Query("SELECT id, name, description, created_at, updated_at, author_id FROM tasks WHERE author_id = $1 LIMIT $2 OFFSET $3", authorID.String(), limit, (page-1)*limit)
	if err != nil {
		return []domain.Task{}, 0, fmt.Errorf("%w: %w", domain.ErrGetTasks, err)
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)

	for rows.Next() {
		t := domain.Task{}
		err = rows.Scan(&t.ID, &t.Name, &t.Description, &t.CreatedAt, &t.UpdatedAt, &t.AuthorID)
		if err != nil {
			return []domain.Task{}, 0, fmt.Errorf("get tasks: %w", err)
		}
		tasks = append(tasks, t)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return []domain.Task{}, 0, fmt.Errorf("%w: %w", domain.ErrGetTasks, rowsErr)
	}

	var totalCount int

	err = r.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE author_id = $1", authorID).Scan(&totalCount)
	if err != nil {
		return []domain.Task{}, 0, fmt.Errorf("%w: %w", domain.ErrGetTasks, err)
	}

	return tasks, totalCount, nil
}
