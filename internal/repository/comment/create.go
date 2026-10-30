package comment

import (
	"TaskFlow/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *commentRepository) CreateComment(
	ctx context.Context,
	ID, taskID, authorID uuid.UUID,
	text string,
) (domain.TaskComment, error) {
	var comment domain.TaskComment

	err := r.db.QueryRowContext(ctx, `
					WITH inserted AS (
					    INSERT INTO task_comments (id, task_id, author_id, text, created_at)
					    SELECT $1, $2, $3, $4, $5
						WHERE EXISTS (SELECT 1 FROM tasks WHERE id = $2)
					    RETURNING id, task_id, author_id, text, created_at
					)
					SELECT i.id, i.task_id, i.author_id, COALESCE(p.username, '') AS author_name, i.text, i.created_at
					FROM inserted i
					LEFT JOIN profiles p ON p.user_id = i.author_id`, ID, taskID, authorID, text, time.Now()).
		Scan(&comment.ID, &comment.TaskID, &comment.AuthorID, &comment.AuthorName, &comment.Text, &comment.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TaskComment{}, domain.ErrTaskNotFound
		}

		return domain.TaskComment{}, fmt.Errorf("insert comment %s: %w", ID, err)
	}
	return comment, nil
}
