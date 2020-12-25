package handlers

import (
	"strconv"

	"github.com/dgrijalva/jwt-go"
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

// validatePassword checks the password against the salted hash.
func validatePassword(hash string, pwd string) bool {
	byteHash := []byte(hash)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(pwd))
	if err != nil {
		return false
	}
	return true
}

// getUserIDFromClaim gets the userID from the JWT token.
func getUserIDFromClaim(c echo.Context) (int64, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	claimUID := claims["user_id"]
	userID, err := strconv.ParseInt(claimUID.(string), 10, 64)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
