package store

import "database/sql"

type Repository struct {
	DbHandle *sql.DB
}
