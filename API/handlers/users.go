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
	limit, offset := getLimitAndOffset(c)

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
	id, _ := strconv.Atoi(c.Param("id"))
	u, err := models.UserByUserID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't process input user")
	}

	if err := c.Bind(u); err != nil {
		return err
	}

	u.Email = strings.ToLower(u.Email)

	err = u.Update(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save user "+err.Error())
	}

	return c.JSON(http.StatusCreated, u)
}
