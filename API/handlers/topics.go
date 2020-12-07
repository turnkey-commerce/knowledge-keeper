package handlers

import (
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"github.com/turnkey-commerce/knowledge-keeper/models"
	nullable "gopkg.in/guregu/null.v4"
)

// GetRecentTopicsPaginated returns the recently added categories by limit/offset.
func (h *Handler) GetRecentTopicsPaginated(c echo.Context) error {
	limit, offset := getLimitAndOffset(c)

	topics, err := models.GetRecentPaginatedTopics(h.DB, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find recent topics: "+err.Error())
	}

	return c.JSON(http.StatusOK, topics)
}

// GetTopicByTitle returns the topic by title.
func (h *Handler) GetTopicByTitle(c echo.Context) error {
	title := c.Param("title")
	topics, err := models.TopicsByTitle(h.DB, title)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find topic "+title)
	}

	return c.JSON(http.StatusOK, topics)
}

// SaveTopic saves the topic to the database.
func (h *Handler) SaveTopic(c echo.Context) error {
	t := &models.Topic{}
	if err := c.Bind(t); err != nil {
		return err
	}

	// TODO: get userId from token
	userID := 4
	t.CreatedBy = int64(userID)
	t.UpdatedBy = nullable.Int{}

	err := t.Save(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save topic "+err.Error())
	}

	return c.JSON(http.StatusCreated, t)
}

// UpdateTopic updates the category in the database.
func (h *Handler) UpdateTopic(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	t, err := models.TopicByTopicID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't process input category")
	}

	if err := c.Bind(t); err != nil {
		return err
	}

	// TODO: get userId from token
	userID := 4
	t.UpdatedBy = nullable.IntFrom(int64(userID))

	err = t.Update(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save topic "+err.Error())
	}

	return c.JSON(http.StatusCreated, t)
}
