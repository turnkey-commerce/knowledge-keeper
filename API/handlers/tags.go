package handlers

import (
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"github.com/turnkey-commerce/knowledge-keeper/models"
	"github.com/turnkey-commerce/knowledge-keeper/models/queries"
	nullable "gopkg.in/guregu/null.v4"
)

// GetRecentTagsPaginated returns the recently added tags by limit/offset.
func (h *Handler) GetRecentTagsPaginated(c echo.Context) error {
	limit, offset := getLimitAndOffset(c)

	tags, err := models.GetRecentPaginatedTags(h.DB, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find recent tags: "+err.Error())
	}

	return c.JSON(http.StatusOK, tags)
}

// GetTagByTagID returns the tag by id.
func (h *Handler) GetTagByTagID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	tag, err := models.TagByTagID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't get tag by ID "+c.Param("id"))
	}

	return c.JSON(http.StatusOK, tag)
}

// GetTagTopics returns the topics associated with the tag.
func (h *Handler) GetTagTopics(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	topics, err := queries.TopicsByTagID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find topics")
	}

	return c.JSON(http.StatusOK, topics)
}

// GetTagByName returns the tag by name.
func (h *Handler) GetTagByName(c echo.Context) error {
	name := c.Param("name")
	tags, err := models.TagsByName(h.DB, name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find tag "+name)
	}

	return c.JSON(http.StatusOK, tags)
}

// SaveTag saves the tag to the database.
func (h *Handler) SaveTag(c echo.Context) error {
	t := &models.Tag{}
	if err := c.Bind(t); err != nil {
		return err
	}

	userID, err := getUserIDFromClaim(c)
	if err != nil {
		return err
	}
	t.CreatedBy = userID
	t.UpdatedBy = nullable.Int{}

	err = t.Save(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save tag "+err.Error())
	}

	return c.JSON(http.StatusCreated, t)
}

// UpdateTag updates the tag in the database.
func (h *Handler) UpdateTag(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	t, err := models.TagByTagID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't process input tag")
	}

	if err := c.Bind(t); err != nil {
		return err
	}

	userID, err := getUserIDFromClaim(c)
	if err != nil {
		return err
	}
	t.UpdatedBy = nullable.IntFrom(int64(userID))

	err = t.Update(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save tag "+err.Error())
	}

	return c.JSON(http.StatusCreated, t)
}
