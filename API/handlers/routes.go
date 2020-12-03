package handlers

import (
	echo "github.com/labstack/echo/v4"
)

// GetRoutes routes to the handlers
func (h *Handler) GetRoutes(e *echo.Echo) {
	e.GET("/users/email/:email", h.GetUserByEmail)
	e.POST("/users", h.SaveUser)
	e.PUT("/users/:id", h.UpdateUser)
}
