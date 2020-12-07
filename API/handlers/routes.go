package handlers

import (
	echo "github.com/labstack/echo/v4"
)

// GetRoutes routes to the handlers
func (h *Handler) GetRoutes(e *echo.Echo) {
	// Users
	e.GET("/users", h.GetRecentUsersPaginated)
	e.GET("/users/email/:email", h.GetUserByEmail)
	e.POST("/users", h.SaveUser)
	e.PUT("/users/:id", h.UpdateUser)
	// Categories
	e.GET("/categories", h.GetRecentCategoriesPaginated)
	e.GET("/categories/name/:name", h.GetCategoryByName)
	e.POST("/categories", h.SaveCategory)
	e.PUT("/categories/:id", h.UpdateCategory)
	// Topics
	e.GET("/topics", h.GetRecentTopicsPaginated)
	e.GET("/topics/name/:title", h.GetTopicByTitle)
	e.POST("/topics", h.SaveTopic)
	e.PUT("/topics/:id", h.UpdateTopic)

}
