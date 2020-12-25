package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/xo/dburl"

	"github.com/turnkey-commerce/knowledge-keeper/config"
	"github.com/turnkey-commerce/knowledge-keeper/models"
)

// GetConnectionString returns a DB Connection string
func GetConnectionString(conf config.DatabaseConfigurations) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", conf.DBUser, conf.DBPassword, conf.Server, conf.DBName)
}

// New creates a new SQL DB
func New(connString string) *sql.DB {
	db, err := dburl.Open(connString)
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	return db
}

// Seed the database with an admin user
func Seed(db *sql.DB, email string, password string) error {
	// Check that there's not a duplicate user
	users, _ := models.UsersByEmail(db, email)
	var hash string
	var err error
	if len(users) == 0 {
		hash, err = HashPassword(password)
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

// HashPassword returns the salted hash for a given password.
func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}
