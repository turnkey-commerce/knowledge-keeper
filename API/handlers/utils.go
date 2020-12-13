package handlers

import (
	"strconv"

	echo "github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// getLimitAndOffset will pass
func getLimitAndOffset(c echo.Context) (limit int, offset int) {
	limitIn, _ := strconv.Atoi(c.QueryParam("limit"))
	offsetIn, _ := strconv.Atoi(c.QueryParam("offset"))

	// In case no limit is passed default to 50.
	if limitIn == 0 {
		limitIn = 50
	}
	return limitIn, offsetIn
}

// hashPassword returns the salted hash for a given password.
func hashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

// validatePassword checks the password against the salted hash.
func validatePassword(hash string, pwd string) bool {
	byteHash := []byte(hash)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(pwd))
	if err != nil {
		return false
	}
	return true
}
