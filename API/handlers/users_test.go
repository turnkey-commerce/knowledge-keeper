package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	txdb "github.com/DATA-DOG/go-txdb"
	"github.com/dgrijalva/jwt-go"
	echo "github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/turnkey-commerce/knowledge-keeper/models"
	"github.com/xo/dburl"
)

var (
	email1       = "juke@email.com"
	email1Update = "juke2@email.com"
	email2       = "jack@email.com"

	user1JSON = fmt.Sprintf(`{"email": "%s", 
    "first_name": "Juke",
    "last_name": "Smith",
    "password": "Test-Password",
	"is_admin": false,
	"is_active": true}`, email1)

	user1UpdateJSON = fmt.Sprintf(`{"email": "%s", 
    "first_name": "Jukie",
    "last_name": "Smithy",
	"is_admin": true,
	"is_active": false}`, email1Update)

	user2JSON = fmt.Sprintf(`{"email": "%s", 
    "first_name": "Jack",
    "last_name": "Smith",
    "password": "Test-Password",
	"is_admin": false,
	"is_active": true}`, email2)

	user2UpdateJSON = fmt.Sprintf(`{"email": "%s", 
    "first_name": "Jackie",
    "last_name": "Smithy",
    "password": "Test-Password",
	"is_admin": true}`, email2)

	user1LoginJSON = fmt.Sprintf(`{"email": "%s", 
	"password": "Test-Password"}`, email1)

	user1BadLoginJSON = fmt.Sprintf(`{"email": "%s", 
	"password": "Bogus"}`, email1)

	adminUserLoginJSON = `{"email": "admin@example.com", 
	"password": "change$me"}`
)

func init() {
	u, err := dburl.Parse("postgres://knowledge-keeper:knowledge-keeper@localhost/knowledge-keeper?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	txdb.Register("knowledge", "postgres", u.DSN)
}

// TestRegisterUserCantCreateDuplicate tests the route to create a new user.
// Also validates that a duplicate user can't be created.
func TestRegisterUserAndCantCreateDuplicate(t *testing.T) {
	// Setup
	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	rec := httptest.NewRecorder()
	e := echo.New()
	h := NewHandler(db, "secret")

	// Register the user and make assertions
	if assert.NoError(t, registerUser(db, e, rec, h, user1JSON)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		result := make(map[string]interface{})
		json.Unmarshal([]byte(rec.Body.String()), &result)
		resultEmail := result["email"]

		assert.Equal(t, email1, resultEmail)
	}

	// Assert that there is an error on trying to create the same user again.
	assert.Error(t, registerUser(db, e, rec, h, user1JSON))
}

// TestRegisteredUserCanLogin tests that the registered user can login.
// Also validates that they can't login if wrong password.
func TestRegisteredUserCanLogin(t *testing.T) {
	// Setup
	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	rec := httptest.NewRecorder()
	e := echo.New()
	h := NewHandler(db, "secret")

	// First register the user.
	assert.NoError(t, registerUser(db, e, rec, h, user1JSON))

	// Login Success
	assert.NoError(t, loginUser1(db, e, rec, h, user1LoginJSON))

	// Login Failure
	assert.Error(t, loginUser1(db, e, rec, h, user1BadLoginJSON))
}

// TestCanViewRecentUsers tests that the registered user can view recent users.
func TestCanViewRecentUsers(t *testing.T) {
	// Setup
	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	rec := httptest.NewRecorder()
	e := echo.New()
	h := NewHandler(db, "secret")

	// First register the user.
	assert.NoError(t, registerUser(db, e, rec, h, user1JSON))

	// Login Success
	assert.NoError(t, loginUser1(db, e, rec, h, user1LoginJSON))

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	assert.NoError(t, h.GetRecentUsersPaginated(c))
	var users []models.User
	json.Unmarshal([]byte(rec.Body.String()), &users)
	assert.Equal(t, email1, users[0].Email)
}

// TestCanGetUserByEmail tests that the registered user can get user by email.
func TestCanGetUserByEmail(t *testing.T) {
	// Setup
	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	rec := httptest.NewRecorder()
	e := echo.New()
	h := NewHandler(db, "secret")

	// First register the user.
	assert.NoError(t, registerUser(db, e, rec, h, user1JSON))

	// Login Success
	assert.NoError(t, loginUser1(db, e, rec, h, user1LoginJSON))

	// Setup to get the user by email.
	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/email/{email}", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.SetPath("/users/email/:email")
	c.SetParamNames("email")
	c.SetParamValues(email1)

	assert.NoError(t, h.GetUserByEmail(c))
	var users []models.User
	json.Unmarshal([]byte(rec.Body.String()), &users)
	assert.Equal(t, email1, users[0].Email)
}

// TestCanUpdateOwnUserOnly tests that the registered user can update their
//  allowed properties but not another user's properties.
func TestCanUpdateOwnUserOnly(t *testing.T) {
	// Setup
	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	rec := httptest.NewRecorder()
	e := echo.New()
	h := NewHandler(db, "secret")

	// First register the user.
	assert.NoError(t, registerUser(db, e, rec, h, user1JSON))
	// Get the ID of the registered user
	var user models.User
	json.Unmarshal([]byte(rec.Body.String()), &user)
	id := user.UserID

	// Login the user
	rec = httptest.NewRecorder()
	assert.NoError(t, loginUser1(db, e, rec, h, user1LoginJSON))
	token := getLoginToken(rec)

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users/{id}", strings.NewReader(user1UpdateJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))
	c.Set("user", parseToken(token))

	assert.NoError(t, h.UpdateUser(c))

	var updatedUser models.User
	json.Unmarshal([]byte(rec.Body.String()), &updatedUser)
	assert.Equal(t, email1Update, updatedUser.Email)
	// Check that they can't promote themselves to an admin
	assert.Equal(t, false, updatedUser.IsAdmin)
	// Check that they can't inactivate themselves
	assert.Equal(t, true, updatedUser.IsActive)

	// Check that they can't edit another user
	// Register a second user.
	rec = httptest.NewRecorder()
	assert.NoError(t, registerUser(db, e, rec, h, user2JSON))
	// Get the ID of the second user
	json.Unmarshal([]byte(rec.Body.String()), &user)
	id = user.UserID

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/users/{id}", strings.NewReader(user2UpdateJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))
	c.Set("user", parseToken(token))

	assert.Error(t, h.UpdateUser(c))
}

// TestAdminCanUpdateUsers tests that the admin can update a registered user to change
//   their admin status and make them inactive.
func TestAdminCanUpdateUsers(t *testing.T) {
	// Setup
	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	rec := httptest.NewRecorder()
	e := echo.New()
	h := NewHandler(db, "secret")

	// First register the user.
	assert.NoError(t, registerUser(db, e, rec, h, user1JSON))
	// Get the ID of the registered user
	var user models.User
	json.Unmarshal([]byte(rec.Body.String()), &user)
	id := user.UserID

	// Login the admin user
	rec = httptest.NewRecorder()
	assert.NoError(t, loginUser1(db, e, rec, h, adminUserLoginJSON))
	token := getLoginToken(rec)

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/users/{id}", strings.NewReader(user1UpdateJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))
	c.Set("user", parseToken(token))

	assert.NoError(t, h.UpdateUser(c))

	var updatedUser models.User
	json.Unmarshal([]byte(rec.Body.String()), &updatedUser)
	assert.Equal(t, email1Update, updatedUser.Email)
	// Check that they were promoted to Admin and Inactive
	assert.Equal(t, true, updatedUser.IsAdmin)
	assert.Equal(t, false, updatedUser.IsActive)
}

// Private Functions

func registerUser(db *sql.DB, e *echo.Echo, rec *httptest.ResponseRecorder, h *Handler, json string) error {
	req := httptest.NewRequest(http.MethodPost, "/registration", strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	return h.UserRegistration(c)
}

func loginUser1(db *sql.DB, e *echo.Echo, rec *httptest.ResponseRecorder, h *Handler, json string) error {
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	return h.UserLogin(c)
}

func getLoginToken(rec *httptest.ResponseRecorder) string {
	result := make(map[string]string)
	json.Unmarshal([]byte(rec.Body.String()), &result)
	return string(result["token"])
}

func parseToken(tokenStr string) *jwt.Token {
	token, _ := jwt.Parse(tokenStr, nil)
	return token
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
