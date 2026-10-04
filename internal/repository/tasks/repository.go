package tasks

import "database/sql"

type tasksRepository struct {
	db *sql.DB
}

func NewTasksRepository(db *sql.DB) *tasksRepository {
	return &tasksRepository{db: db}
}
