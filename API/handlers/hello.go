package handlers

import (
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

// GetHello returns hello world.
func (h *Handler) GetHello(c echo.Context) error {
	fmt.Println(h.DB.Stats())
	return c.String(http.StatusOK, "Hello World")
}
