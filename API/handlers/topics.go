package handlers

import (
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"github.com/turnkey-commerce/knowledge-keeper/models"
	"github.com/turnkey-commerce/knowledge-keeper/models/queries"
	nullable "gopkg.in/guregu/null.v4"
)

// GetRecentTopicsPaginated returns the recently added topics by limit/offset.
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

// GetTopicTags returns the tags associated with the topic.
func (h *Handler) GetTopicTags(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	tags, err := queries.TagsByTopicID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find tags")
	}

	return c.JSON(http.StatusOK, tags)
}

// GetTopicNotes returns the notes associated with the topic.
func (h *Handler) GetTopicNotes(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	tags, err := queries.NotesByTopicID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find notes")
	}

	return c.JSON(http.StatusOK, tags)
}

// GetTopicMedia returns the media associated with the topic.
func (h *Handler) GetTopicMedia(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	tags, err := queries.MediaByTopicID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find media")
	}

	return c.JSON(http.StatusOK, tags)
}

// GetRelatedTopics returns the related topics with the topic.
func (h *Handler) GetRelatedTopics(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	topics, err := queries.RelatedTopicsByTopicId(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find related topics")
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
	userID := 1
	t.CreatedBy = int64(userID)
	t.UpdatedBy = nullable.Int{}

	err := t.Save(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save topic "+err.Error())
	}

	return c.JSON(http.StatusCreated, t)
}

// UpdateTopic updates the topic in the database.
func (h *Handler) UpdateTopic(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	t, err := models.TopicByTopicID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't process input topic")
	}

	if err := c.Bind(t); err != nil {
		return err
	}

	// TODO: get userId from token
	userID := 1
	t.UpdatedBy = nullable.IntFrom(int64(userID))

	err = t.Update(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save topic "+err.Error())
	}

	return c.JSON(http.StatusCreated, t)
}

// AddTagToTopic saves a tag to the topic.
func (h *Handler) AddTagToTopic(c echo.Context) error {
	t := &models.TopicsTag{}
	if err := c.Bind(t); err != nil {
		return err
	}

	err := t.Insert(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't add topic to tag "+err.Error())
	}

	return c.JSON(http.StatusCreated, t)
}

// AddRelatedTopic saves a related topic.
func (h *Handler) AddRelatedTopic(c echo.Context) error {
	t := &models.RelatedTopic{}
	if err := c.Bind(t); err != nil {
		return err
	}

	err := t.Insert(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't add related topic "+err.Error())
	}

	return c.JSON(http.StatusCreated, t)
}
