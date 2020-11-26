package database_test

import (
	"database/sql"
	"log"
	"testing"

	txdb "github.com/DATA-DOG/go-txdb"
	"github.com/turnkey-commerce/knowledge-keeper/models"
	"github.com/xo/dburl"
)

func init() {
	u, err := dburl.Parse("postgres://knowledge-keeper:knowledge-keeper@localhost/knowledge-keeper?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	txdb.Register("knowledge", "postgres", u.DSN)
}

// TestCreateUserAndCategories tests creating a user and adding
func TestCreateUserAndCategories(t *testing.T) {
	db, err := sql.Open("knowledge", "identifier")

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Email:     "test@test.com",
		FirstName: sql.NullString{String: "Jack"},
		LastName:  sql.NullString{String: "Test"},
	}

	err = user.Save(db)
	if err != nil {
		log.Fatal(err)
	}

	category := models.Category{
		Name:        "Special",
		Description: sql.NullString{String: "Special Description"},
		CreatedBy:   user.UserID,
	}

	err = category.Save(db)
	if err != nil {
		log.Fatal(err)
	}
}
