package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/xo/dburl"

	"github.com/turnkey-commerce/knowledge-keeper/handlers"
	"github.com/turnkey-commerce/knowledge-keeper/models"
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

func Seed(db *sql.DB, email string, password string) error {
	// Check that there's not a duplicate user
	users, _ := models.UsersByEmail(db, email)
	var hash string
	var err error
	if len(users) == 0 {
		hash, err = handlers.HashPassword(password)
		u := models.User{
			Email:     email,
			FirstName: "Admin",
			LastName:  "",
			IsAdmin:   true,
			IsActive:  true,
			Hash:      hash,
		}

		err = u.Save(db)
		if err != nil {
			return err
		}
	}

	return nil
}
