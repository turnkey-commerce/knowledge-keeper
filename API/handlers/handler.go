package handlers

import (
	"database/sql"
)

// Handler carries the DB and secret key references
type Handler struct {
	DB     *sql.DB
	Secret string
}

// NewHandler returns a Handler
func NewHandler(db *sql.DB, secret string) *Handler {
	return &Handler{
		DB:     db,
		Secret: secret,
	}
}
