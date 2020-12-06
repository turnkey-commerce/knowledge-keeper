package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"github.com/turnkey-commerce/knowledge-keeper/models"
)

// GetRecentCategoriesPaginated returns the recently added categories by limit/offset.
func (h *Handler) GetRecentCategoriesPaginated(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	// In case no limit is passed default to 50.
	if limit == 0 {
		limit = 50
	}

	categories, err := models.GetRecentPaginatedCategorys(h.DB, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find recent categories: "+err.Error())
	}

	return c.JSON(http.StatusOK, categories)
}

// GetCategoryByName returns the category by name.
func (h *Handler) GetCategoryByName(c echo.Context) error {
	name := c.Param("name")
	categories, err := models.CategoriesByName(h.DB, name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find category "+name)
	}

	return c.JSON(http.StatusOK, categories)
}

// SaveCategory saves the category to the database.
func (h *Handler) SaveCategory(c echo.Context) error {
	cat := &models.Category{}
	if err := c.Bind(cat); err != nil {
		return err
	}

	// TODO: get userId from token
	userID := 4
	cat.CreatedBy = int64(userID)
	cat.UpdatedBy = sql.NullInt64{}

	err := cat.Save(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save category "+err.Error())
	}

	return c.JSON(http.StatusCreated, c)
}

// UpdateCategory updates the category in the database.
func (h *Handler) UpdateCategory(c echo.Context) error {
	cat := &models.Category{}
	if err := c.Bind(cat); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	cat.CategoryID = int64(id)
	// TODO: get userId from token
	userID := 4
	cat.UpdatedBy = sql.NullInt64{Int64: int64(userID)}

	err := cat.Upsert(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save category "+err.Error())
	}

	return c.JSON(http.StatusCreated, cat)
}
