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
	e.GET("/topics/:id/tags", h.GetTopicTags)
	e.GET("/topics/title/:title", h.GetTopicByTitle)
	e.POST("/topics", h.SaveTopic)
	e.PUT("/topics/:id", h.UpdateTopic)
	e.POST("/topics/tag", h.AddTagToTopic)
	// Tags
	e.GET("/tags", h.GetRecentTagsPaginated)
	e.GET("/tags/:id/topics", h.GetTagTopics)
	e.GET("/tags/name/:name", h.GetTagByName)
	e.POST("/tags", h.SaveTag)
	e.PUT("/tags/:id", h.UpdateTag)

}
