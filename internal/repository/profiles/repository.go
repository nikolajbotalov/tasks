package profiles

import "database/sql"

type profilesRepository struct {
	db *sql.DB
}

func NewProfilesRepository(db *sql.DB) *profilesRepository {
	return &profilesRepository{db: db}
}
