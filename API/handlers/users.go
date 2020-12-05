package handlers

import (
	"net/http"
	"strconv"
	"strings"

	echo "github.com/labstack/echo/v4"
	"github.com/turnkey-commerce/knowledge-keeper/models"
)

// GetRecentUsersPaginated returns the recently added users by limit/offset.
func (h *Handler) GetRecentUsersPaginated(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	// In case no limit is passed default to 50.
	if limit == 0 {
		limit = 50
	}

	users, err := models.GetRecentPaginatedUsers(h.DB, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find recent users: "+err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

// GetUserByEmail returns the user by email address.
func (h *Handler) GetUserByEmail(c echo.Context) error {
	email := c.Param("email")
	users, err := models.UsersByEmail(h.DB, email)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find user "+email)
	}

	return c.JSON(http.StatusOK, users)
}

// SaveUser saves the user to the database.
func (h *Handler) SaveUser(c echo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return err
	}
	u.Email = strings.ToLower(u.Email)

	err := u.Save(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save user "+err.Error())
	}

	return c.JSON(http.StatusCreated, u)
}

// UpdateUser updates the user in the database.
func (h *Handler) UpdateUser(c echo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	u.Email = strings.ToLower(u.Email)
	u.UserID = int64(id)

	err := u.Upsert(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save user "+err.Error())
	}

	return c.JSON(http.StatusCreated, u)
}
