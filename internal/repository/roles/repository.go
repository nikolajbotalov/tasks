package roles

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRolesRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}
