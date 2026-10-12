package permissions

import "database/sql"

type permissionsRepository struct {
	db *sql.DB
}

func NewPermissionsRepository(db *sql.DB) *permissionsRepository {
	return &permissionsRepository{db: db}
}
