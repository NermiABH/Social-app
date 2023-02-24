package store

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}
