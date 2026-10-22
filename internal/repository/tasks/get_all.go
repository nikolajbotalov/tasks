package tasks

import (
	"TaskFlow/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

func (r *tasksRepository) GetTaskList(
	authorID uuid.UUID,
	assigneeID *uuid.UUID,
	status domain.TaskStatus,
	priority domain.TaskPriority,
	search string,
	page, limit int,
) ([]domain.Task, int, error) {
	escapedSearch := escapeLike(search)
	rows, err := r.db.Query(`
		SELECT id, name, description, created_at, updated_at, author_id, assignee_id, status, priority 
		FROM tasks
		WHERE author_id = $1 
		  AND ($2::uuid IS NULL OR assignee_id = $2::uuid) 
		  AND ($3 = '' OR status = $3) 
		  AND ($4 = '' OR priority = $4)
		  AND (
		      $5 = '' 
		      OR name ILIKE '%' || $5 || '%' ESCAPE '\'
		      OR description ILIKE '%' || $5 || '%' ESCAPE '\'
		      )
		LIMIT $6 OFFSET $7`,
		authorID, assigneeID, status, priority, escapedSearch, limit, (page-1)*limit)
	if err != nil {
		return []domain.Task{}, 0, fmt.Errorf("query task list: %w", err)
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)

	for rows.Next() {
		t := domain.Task{}
		err = rows.Scan(&t.ID, &t.Name, &t.Description, &t.CreatedAt, &t.UpdatedAt, &t.AuthorID, &t.AssigneeID, &t.Status, &t.Priority)
		if err != nil {
			return []domain.Task{}, 0, fmt.Errorf("scan task %s: %w", t.ID, err)
		}
		tasks = append(tasks, t)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return []domain.Task{}, 0, fmt.Errorf("iterate task list: %w", rowsErr)
	}

	var totalCount int

	err = r.db.QueryRow(`
		SELECT COUNT(*) FROM tasks
		WHERE author_id = $1
		  AND ($2::uuid IS NULL OR assignee_id = $2::uuid) 
		  AND ($3 = '' OR status = $3) 
		  AND ($4 = '' OR priority = $4)
		  AND (
		      $5 = '' 
		      OR name ILIKE '%' || $5 || '%' ESCAPE '\'
		      OR description ILIKE '%' || $5 || '%' ESCAPE '\'
		      )`,
		authorID, assigneeID, status, priority, escapedSearch).
		Scan(&totalCount)
	if err != nil {
		return []domain.Task{}, 0, fmt.Errorf("count task: %w", err)
	}

	return tasks, totalCount, nil
}
