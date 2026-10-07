package comment

import "database/sql"

type commentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *commentRepository {
	return &commentRepository{
		db: db,
	}
}
