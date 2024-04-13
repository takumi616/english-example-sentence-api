package store

import "database/sql"

type repository struct {
	DbHandle *sql.DB
}
