package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func (h *Handler) UserLogin(c echo.Context) error {
	u := &models.UserAuth{}
	if err := c.Bind(u); err != nil {
		return err
	}
	users, err := models.UsersByEmail(h.DB, u.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find user "+u.Email)
	}
	user := users[0]
	if !validatePassword(user.Hash, u.Password) {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.UserID
	claims["email"] = user.Email
	claims["name"] = user.FirstName + " " + user.LastName
	claims["admin"] = user.IsAdmin
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	// TODO: Get secret from config.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

// SaveUser saves the user to the database.
func (h *Handler) SaveUser(c echo.Context) error {
	u := &models.UserAuth{}
	if err := c.Bind(u); err != nil {
		return err
	}
	u.Email = strings.ToLower(u.Email)
	hash, err := hashPassword(u.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't process user hash")
	}
	u.Hash = hash
	// Clear password so it's not returned as part of the extended struct.
	u.Password = ""
	err = u.Save(h.DB)
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
