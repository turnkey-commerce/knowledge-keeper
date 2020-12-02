package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/xo/dburl"
)

// New creates a new SQL DB
func New() *sql.DB {
	db, err := dburl.Open("postgres://knowledge-keeper:knowledge-keeper@localhost/knowledge-keeper?sslmode=disable")
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	return db
}
