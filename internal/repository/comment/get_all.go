package comment

import (
	"TaskFlow/internal/domain"
	"fmt"

	"github.com/google/uuid"
)

func (r *commentRepository) GetCommentList(taskID uuid.UUID, page, limit int) ([]domain.TaskComment, int, error) {
	rows, err := r.db.Query(
		`SELECT c.id, c.task_id, c.author_id, u.username AS author_name, c.text, c.created_at
				FROM task_comments c JOIN users u ON u.id = c.author_id
				WHERE c.task_id = $1 ORDER BY c.created_at LIMIT $2 OFFSET $3`, taskID, limit, (page-1)*limit)
	if err != nil {
		return []domain.TaskComment{}, 0, fmt.Errorf("%w: %w", domain.ErrGetComments, err)
	}

	defer rows.Close()

	commentList := make([]domain.TaskComment, 0)
	for rows.Next() {
		comment := domain.TaskComment{}
		err = rows.Scan(&comment.ID, &comment.TaskID, &comment.AuthorID, &comment.AuthorName, &comment.Text, &comment.CreatedAt)
		if err != nil {
			return []domain.TaskComment{}, 0, fmt.Errorf("%w: %w", domain.ErrGetComments, err)
		}
		commentList = append(commentList, comment)
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		return []domain.TaskComment{}, 0, fmt.Errorf("%w: %w", domain.ErrGetComments, rowsErr)
	}

	var total int
	err = r.db.QueryRow("SELECT COUNT(*) FROM task_comments WHERE task_id = $1", taskID).Scan(&total)
	if err != nil {
		return []domain.TaskComment{}, 0, fmt.Errorf("%w: %w", domain.ErrGetComments, err)
	}

	return commentList, total, nil
}
