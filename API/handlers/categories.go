package handlers

import (
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"github.com/turnkey-commerce/knowledge-keeper/models"
	nullable "gopkg.in/guregu/null.v4"
)

// GetRecentCategoriesPaginated returns the recently added categories by limit/offset.
func (h *Handler) GetRecentCategoriesPaginated(c echo.Context) error {
	limit, offset := getLimitAndOffset(c)

	categories, err := models.GetRecentPaginatedCategorys(h.DB, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find recent categories: "+err.Error())
	}

	return c.JSON(http.StatusOK, categories)
}

// GetCategoryByCategoryID returns the category by id.
func (h *Handler) GetCategoryByCategoryID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	category, err := models.CategoryByCategoryID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't get category by ID "+c.Param("id"))
	}

	return c.JSON(http.StatusOK, category)
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

	userID, err := getUserIDFromClaim(c)
	if err != nil {
		return err
	}
	cat.CreatedBy = userID
	cat.UpdatedBy = nullable.Int{}

	err = cat.Save(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save category "+err.Error())
	}

	return c.JSON(http.StatusCreated, cat)
}

// UpdateCategory updates the category in the database.
func (h *Handler) UpdateCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	cat, err := models.CategoryByCategoryID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't process input category")
	}

	if err := c.Bind(cat); err != nil {
		return err
	}

	userID, err := getUserIDFromClaim(c)
	if err != nil {
		return err
	}
	cat.UpdatedBy = nullable.IntFrom(int64(userID))

	err = cat.Update(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save category "+err.Error())
	}

	return c.JSON(http.StatusCreated, cat)
}
