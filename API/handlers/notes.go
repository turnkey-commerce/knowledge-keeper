package handlers

import (
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"github.com/turnkey-commerce/knowledge-keeper/models"
	nullable "gopkg.in/guregu/null.v4"
)

// GetRecentNotesPaginated returns the recently added notes by limit/offset.
func (h *Handler) GetRecentNotesPaginated(c echo.Context) error {
	limit, offset := getLimitAndOffset(c)

	notes, err := models.GetRecentPaginatedNotes(h.DB, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find recent notes: "+err.Error())
	}

	return c.JSON(http.StatusOK, notes)
}

// GetNoteByNoteID returns the note by id.
func (h *Handler) GetNoteByNoteID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	note, err := models.NoteByNoteID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't get note by ID "+c.Param("id"))
	}

	return c.JSON(http.StatusOK, note)
}

// GetNoteByTitle returns the note by name.
func (h *Handler) GetNoteByTitle(c echo.Context) error {
	title := c.Param("title")
	notes, err := models.NotesByTitle(h.DB, title)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find note "+title)
	}

	return c.JSON(http.StatusOK, notes)
}

// SaveNote saves the note to the database.
func (h *Handler) SaveNote(c echo.Context) error {
	n := &models.Note{}
	if err := c.Bind(n); err != nil {
		return err
	}

	userID, err := getUserIDFromClaim(c)
	if err != nil {
		return err
	}
	n.CreatedBy = userID
	n.UpdatedBy = nullable.Int{}

	err = n.Save(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save note "+err.Error())
	}

	return c.JSON(http.StatusCreated, n)
}

// UpdateNote updates the note in the database.
func (h *Handler) UpdateNote(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	n, err := models.NoteByNoteID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't process input note")
	}

	if err := c.Bind(n); err != nil {
		return err
	}

	userID, err := getUserIDFromClaim(c)
	if err != nil {
		return err
	}
	n.UpdatedBy = nullable.IntFrom(int64(userID))

	err = n.Update(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save note "+err.Error())
	}

	return c.JSON(http.StatusCreated, n)
}
