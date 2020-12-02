package handlers

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Handler carries the DB reference
type Handler struct {
	DB *sql.DB
}

// NewHandler returns a Handler
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		DB: db,
	}
}
