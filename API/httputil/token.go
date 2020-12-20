package httputil

import (
	"github.com/dgrijalva/jwt-go"
	echo "github.com/labstack/echo/v4"
)

// IsAdmin returns whether the user is an admin
func IsAdmin(c echo.Context) bool {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	isAdmin := claims["admin"].(bool)
	return isAdmin
}

// IsUserEditingSelf returns whether the user is editing their user
func IsUserEditingSelf(c echo.Context, userEmail string) bool {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	return email == userEmail
}
