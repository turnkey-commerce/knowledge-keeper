package handlers

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	echo "github.com/labstack/echo/v4"
	"github.com/turnkey-commerce/knowledge-keeper/httputil"
	"github.com/turnkey-commerce/knowledge-keeper/models"
)

// GetRecentUsersPaginated godoc
// @Summary Get Recent Users Paginated
// @Description Gets Recent Users with optional Pagination
// @Tags users
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit returned per page" default(50)
// @Param offset query int false "Offset from start row" default(0)
// @Success 200 {array} models.UserRegistration
// @Failure 404 {object} httputil.HTTPError "Usrs Not Found"
// @Failure 400 {object} httputil.HTTPError "Bad Request"
// @Failure 500 {object} httputil.HTTPError "Bad Request"
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Router /users [get]
func (h *Handler) GetRecentUsersPaginated(c echo.Context) error {
	limit, offset := getLimitAndOffset(c)

	users, err := models.GetRecentPaginatedUsers(h.DB, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find recent users: "+err.Error())
	}

	clearHash(users)

	return c.JSON(http.StatusOK, users)
}

// GetUserByEmail godoc
// @Summary Get user by email
// @Description Gets user by their email address
// @Tags users
// @Accept  json
// @Produce  json
// @Param email path string true "email address of user"
// @Success 200 {array} models.UserRegistration
// @Failure 404 {object} httputil.HTTPError "User Not Found"
// @Failure 400 {object} httputil.HTTPError "Bad Request"
// @Failure 500 {object} httputil.HTTPError "Bad Request"
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Router /users/email/{email} [get]
func (h *Handler) GetUserByEmail(c echo.Context) error {
	email, err := url.QueryUnescape(c.Param("email"))
	if err != nil {
		return err
	}
	users, err := models.UsersByEmail(h.DB, email)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find user "+email)
	}

	clearHash(users)

	return c.JSON(http.StatusOK, users)
}

// UserLogin godoc
// @Summary User Login
// @Description Logs in the user
// @Tags login
// @Accept  json
// @Produce  json
// @Param userLogin body models.UserLogin true "Login"
// @Success 200 {object} models.Token
// @Failure 401 {object} httputil.HTTPError "Unauthorized"
// @Failure 500 {object} httputil.HTTPError "Bad Input"
// @Router /login [post]
func (h *Handler) UserLogin(c echo.Context) error {
	u := &models.UserLogin{}
	if err := c.Bind(u); err != nil {
		return err
	}
	users, err := models.UsersByEmail(h.DB, u.Email)
	if err != nil || len(users) < 1 {
		return echo.ErrUnauthorized
	}
	user := users[0]
	if !validatePassword(user.Hash, u.Password) {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = strconv.FormatInt(user.UserID, 10)
	claims["email"] = user.Email
	claims["name"] = user.FirstName + " " + user.LastName
	claims["admin"] = user.IsAdmin
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(h.Secret))
	if err != nil {
		return err
	}

	tk := models.Token{Token: t}

	return c.JSON(http.StatusOK, tk)
}

// UserRegistration godoc
// @Summary Create new user
// @Description Creates new user for registration
// @Tags register
// @Accept  json
// @Produce  json
// @Param UserRegistration body models.UserRegistration true "Register"
// @Success 201 {object} models.UserRegistration
// @Failure 400 {object} httputil.HTTPError "Bad Request"
// @Failure 401 {object} httputil.HTTPError "Unauthorized"
// @Failure 500 {object} httputil.HTTPError "Bad Input"
// @Router /register [post]
func (h *Handler) UserRegistration(c echo.Context) error {
	u := &models.UserRegistration{}
	if err := c.Bind(u); err != nil {
		return err
	}

	u.Email = strings.ToLower(u.Email)
	// Check that there's not a duplicate user
	users, _ := models.UsersByEmail(h.DB, u.Email)
	if len(users) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "User already exists")
	}

	hash, err := hashPassword(u.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't process user hash")
	}
	u.Hash = hash

	// Registering users aren't admins.
	u.IsAdmin = false

	err = u.Save(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't save user "+err.Error())
	}

	// Clear password and hash so they're not returned as part of the extended struct.
	u.Password = ""
	u.Hash = ""

	return c.JSON(http.StatusCreated, u)
}

// UpdateUser godoc
// @Summary Updates a user
// @Description Updates a user
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "id of user"
// @Param UserUpdate body models.UserUpdate true "Update"
// @Success 201 {object} models.UserRegistration
// @Failure 400 {object} httputil.HTTPError "Bad Request"
// @Failure 401 {object} httputil.HTTPError "Unauthorized"
// @Failure 500 {object} httputil.HTTPError "Bad Input"
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	existingUser, err := models.UserByUserID(h.DB, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't process input user")
	}

	isAdmin := httputil.IsAdmin(c)
	isExistingUser := httputil.IsUserEditingSelf(c, existingUser.Email)

	u := &models.UserUpdate{}
	if err := c.Bind(u); err != nil {
		return err
	}

	if isAdmin {
		// Only Admin's can change this property
		existingUser.IsAdmin = u.IsAdmin
	}

	if isAdmin || isExistingUser {
		// Only admins or the same user can edit.
		hash, err := hashPassword(u.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Can't process user hash")
		}
		existingUser.Hash = hash
		existingUser.Email = strings.ToLower(u.Email)
		existingUser.FirstName = u.FirstName
		existingUser.LastName = u.LastName
	}

	err = existingUser.Update(h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't update user "+err.Error())
	}

	// Clear password the hash so it's not returned.
	existingUser.Hash = ""

	return c.JSON(http.StatusOK, existingUser)
}

func clearHash(users []*models.User) {
	for _, user := range users {
		user.Hash = ""
	}
}
