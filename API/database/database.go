package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/xo/dburl"
)

// New creates a new SQL DB
func New(dbname string, dbserver string, dbuser string, dbpassword string) *sql.DB {
	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbuser, dbpassword, dbserver, dbname)
	db, err := dburl.Open(connString)
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	return db
}
