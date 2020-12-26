package handlers

import (
	"net/http"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"github.com/turnkey-commerce/knowledge-keeper/models"
	nullable "gopkg.in/guregu/null.v4"
)

// GetRecentCategoriesPaginated godoc
// @Summary Get Recent Categories Paginated
// @Description Gets Recent Categories with optional Pagination
// @Tags categories
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit returned per page" default(50)
// @Param offset query int false "Offset from start row" default(0)
// @Success 200 {array} models.CategoryDTO
// @Failure 404 {object} httputil.HTTPError "Categories Not Found"
// @Failure 400 {object} httputil.HTTPError "Bad Request"
// @Failure 500 {object} httputil.HTTPError "Bad Request"
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Router /categories [get]
func (h *Handler) GetRecentCategoriesPaginated(c echo.Context) error {
	limit, offset := getLimitAndOffset(c)

	categories, err := models.GetRecentPaginatedCategorys(h.DB, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find recent categories: "+err.Error())
	}

	return c.JSON(http.StatusOK, categories)
}

// GetCategoryByCategoryID godoc
// @Summary Get Category By ID
// @Description Get Category By ID
// @Tags categories
// @Accept  json
// @Produce  json
// @Param id path int true "id of category"
// @Success 200 {array} models.CategoryDTO
// @Failure 404 {object} httputil.HTTPError "Category Not Found"
// @Failure 400 {object} httputil.HTTPError "Bad Request"
// @Failure 500 {object} httputil.HTTPError "Bad Request"
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Router /categories/{id} [get]
func (h *Handler) GetCategoryByCategoryID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	category, err := models.CategoryByCategoryID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find category by ID "+c.Param("id"))
	}

	return c.JSON(http.StatusOK, category)
}

// GetCategoryByName godoc
// @Summary Get Category By Name
// @Description Get Category By Name
// @Tags categories
// @Accept  json
// @Produce  json
// @Param name path string true "name of category"
// @Success 200 {array} models.CategoryDTO
// @Failure 404 {object} httputil.HTTPError "Category Not Found"
// @Failure 400 {object} httputil.HTTPError "Bad Request"
// @Failure 500 {object} httputil.HTTPError "Bad Request"
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Router /categories/name/{name} [get]
func (h *Handler) GetCategoryByName(c echo.Context) error {
	name := c.Param("name")
	categories, err := models.CategoriesByName(h.DB, name)
	if err != nil || len(categories) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find category "+name)
	}

	return c.JSON(http.StatusOK, categories)
}

// SaveCategory godoc
// @Summary Create new category
// @Description Creates new category
// @Tags categories
// @Accept  json
// @Produce  json
// @Param Category body models.CategoryDTO true "Create Category"
// @Success 201 {object} models.CategoryDTO
// @Failure 400 {object} httputil.HTTPError "Bad Request"
// @Failure 401 {object} httputil.HTTPError "Unauthorized"
// @Failure 500 {object} httputil.HTTPError "Bad Input"
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Router /categories [post]
func (h *Handler) SaveCategory(c echo.Context) error {
	cat := &models.CategoryDTO{}
	if err := c.Bind(cat); err != nil {
		return err
	}

	userID, err := getUserIDFromClaim(c)
	if err != nil {
		return err
	}

	category := models.Category{
		Name:        cat.Name,
		Description: cat.Description,
	}
	category.CreatedBy = userID
	category.UpdatedBy = nullable.Int{}

	err = category.Save(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save category "+err.Error())
	}

	return c.JSON(http.StatusCreated, category)
}

// UpdateCategory godoc
// @Summary Update existing category
// @Description Updates existing category
// @Tags categories
// @Accept  json
// @Produce  json
// @Param id path int true "id of category"
// @Param UpdateCategory body models.CategoryDTO true "Update Category"
// @Success 201 {object} models.CategoryDTO
// @Failure 400 {object} httputil.HTTPError "Bad Request"
// @Failure 401 {object} httputil.HTTPError "Unauthorized"
// @Failure 500 {object} httputil.HTTPError "Bad Input"
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Router /categories/{id} [put]
func (h *Handler) UpdateCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	cat, err := models.CategoryByCategoryID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't process input category")
	}

	updateCat := &models.CategoryDTO{}
	if err := c.Bind(updateCat); err != nil {
		return err
	}

	cat.Name = updateCat.Name
	cat.Description = updateCat.Description

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
