package handlers

import (
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"github.com/turnkey-commerce/knowledge-keeper/models"
	nullable "gopkg.in/guregu/null.v4"
)

// GetRecentMediaPaginated returns the recently added media by limit/offset.
func (h *Handler) GetRecentMediaPaginated(c echo.Context) error {
	limit, offset := getLimitAndOffset(c)

	media, err := models.GetRecentPaginatedMedias(h.DB, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find recent media: "+err.Error())
	}

	return c.JSON(http.StatusOK, media)
}

// GetMediaByMediaID returns the media by id.
func (h *Handler) GetMediaByMediaID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	media, err := models.MediaByMediaID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't get media by ID "+c.Param("id"))
	}

	return c.JSON(http.StatusOK, media)
}

// GetMediaByTitle returns the media by name.
func (h *Handler) GetMediaByTitle(c echo.Context) error {
	title := c.Param("title")
	media, err := models.MediaByTitle(h.DB, title)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find media "+title)
	}

	return c.JSON(http.StatusOK, media)
}

// SaveMedia saves the media to the database.
func (h *Handler) SaveMedia(c echo.Context) error {
	m := &models.Media{}
	if err := c.Bind(m); err != nil {
		return err
	}

	userID, err := getUserIDFromClaim(c)
	if err != nil {
		return err
	}
	m.CreatedBy = userID
	m.UpdatedBy = nullable.Int{}

	err = m.Save(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save media "+err.Error())
	}

	return c.JSON(http.StatusCreated, m)
}

// UpdateMedia updates the media in the database.
func (h *Handler) UpdateMedia(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	m, err := models.MediaByMediaID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't process input media")
	}

	if err := c.Bind(m); err != nil {
		return err
	}

	userID, err := getUserIDFromClaim(c)
	if err != nil {
		return err
	}
	m.UpdatedBy = nullable.IntFrom(int64(userID))

	err = m.Update(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save media "+err.Error())
	}

	return c.JSON(http.StatusCreated, m)
}
