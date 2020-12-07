package handlers

import (
	"strconv"

	echo "github.com/labstack/echo/v4"
)

func getLimitAndOffset(c echo.Context) (limit int, offset int) {
	limitIn, _ := strconv.Atoi(c.QueryParam("limit"))
	offsetIn, _ := strconv.Atoi(c.QueryParam("offset"))

	// In case no limit is passed default to 50.
	if limitIn == 0 {
		limitIn = 50
	}
	return limitIn, offsetIn
}
