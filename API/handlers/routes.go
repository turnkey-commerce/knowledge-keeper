package handlers

import (
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// GetRoutes routes to the handlers
func (h *Handler) GetRoutes(e *echo.Echo) {

	// Register
	e.POST("/register", h.UserRegistration)

	// Login
	e.POST("/login", h.UserLogin)

	// Users Restricted
	u := e.Group("/users")
	u.Use(middleware.JWT([]byte(h.Secret)))

	u.PUT("/:id", h.UpdateUser)
	u.GET("", h.GetRecentUsersPaginated)
	u.GET("/email/:email", h.GetUserByEmail)

	// Categories
	c := e.Group("/categories")
	c.Use(middleware.JWT([]byte(h.Secret)))

	c.GET("", h.GetRecentCategoriesPaginated)
	c.GET("/:id", h.GetCategoryByCategoryID)
	c.GET("/name/:name", h.GetCategoryByName)
	c.POST("", h.SaveCategory)
	c.PUT("/:id", h.UpdateCategory)

	// Topics
	t := e.Group("/topics")
	t.Use(middleware.JWT([]byte(h.Secret)))

	t.GET("", h.GetRecentTopicsPaginated)
	t.GET("/:id", h.GetTopicByTopicID)
	t.GET("/:id/tags", h.GetTopicTags)
	t.GET("/:id/notes", h.GetTopicNotes)
	t.GET("/:id/media", h.GetTopicMedia)
	t.GET("/:id/related", h.GetRelatedTopics)
	t.GET("/title/:title", h.GetTopicByTitle)
	t.POST("", h.SaveTopic)
	t.PUT("/:id", h.UpdateTopic)
	t.POST("/tag", h.AddTagToTopic)
	t.POST("/related", h.AddRelatedTopic)

	// Tags
	tg := e.Group("/tags")
	tg.Use(middleware.JWT([]byte(h.Secret)))

	tg.GET("", h.GetRecentTagsPaginated)
	tg.GET("/:id", h.GetTagByTagID)
	tg.GET("/:id/topics", h.GetTagTopics)
	tg.GET("/name/:name", h.GetTagByName)
	tg.POST("", h.SaveTag)
	tg.PUT("/:id", h.UpdateTag)

	// Notes
	n := e.Group("/notes")
	n.Use(middleware.JWT([]byte(h.Secret)))

	n.GET("", h.GetRecentNotesPaginated)
	n.GET("/:id", h.GetNoteByNoteID)
	n.GET("/title/:title", h.GetNoteByTitle)
	n.POST("", h.SaveNote)
	n.PUT("/:id", h.UpdateNote)

	// Media
	m := e.Group("/media")
	m.Use(middleware.JWT([]byte(h.Secret)))

	m.GET("", h.GetRecentMediaPaginated)
	m.GET("/:id", h.GetMediaByMediaID)
	m.GET("/title/:title", h.GetMediaByTitle)
	m.POST("", h.SaveMedia)
	m.PUT("/:id", h.UpdateMedia)
}
