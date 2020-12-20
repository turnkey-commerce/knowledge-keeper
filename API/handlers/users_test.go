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
	echo "github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/turnkey-commerce/knowledge-keeper/models"
	"github.com/xo/dburl"
)

var (
	email1 = "juke@email.com"
	email2 = "jack@email.com"

	user1JSON = fmt.Sprintf(`{"email": "%s", 
    "first_name": "Juke",
    "last_name": "Smith",
    "password": "Test-Password",
	"is_admin": true}`, email1)

	user2JSON = fmt.Sprintf(`{"email": "%s", 
    "first_name": "Jack",
    "last_name": "Smith",
    "password": "Test-Password",
	"is_admin": true}`, email2)

	user1LoginJSON = fmt.Sprintf(`{"email": "%s", 
	"password": "Test-Password"}`, email1)

	user1BadLoginJSON = fmt.Sprintf(`{"email": "%s", 
	"password": "Bogus"}`, email1)
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
	if assert.NoError(t, registerUser1(db, e, rec, h)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		result := make(map[string]interface{})
		json.Unmarshal([]byte(rec.Body.String()), &result)
		resultEmail := result["email"]

		assert.Equal(t, email1, resultEmail)
	}

	// Assert that there is an error on trying to create the same user again.
	assert.Error(t, registerUser1(db, e, rec, h))
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
	assert.NoError(t, registerUser1(db, e, rec, h))

	// Login Success
	assert.NoError(t, loginUser1(db, e, rec, h, user1LoginJSON))

	// Login Failure
	assert.Error(t, loginUser1(db, e, rec, h, user1BadLoginJSON))
}

// TestRegisteredAdminCanViewUsers tests that the registered admin can view recent users.
func TestRegisteredAdminCanViewUsers(t *testing.T) {
	// Setup
	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	rec := httptest.NewRecorder()
	e := echo.New()
	h := NewHandler(db, "secret")

	// First register the user.
	assert.NoError(t, registerUser1(db, e, rec, h))

	// Login Success
	assert.NoError(t, loginUser1(db, e, rec, h, user1LoginJSON))

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	assert.NoError(t, h.GetRecentUsersPaginated(c))
	var users []models.User
	json.Unmarshal([]byte(rec.Body.String()), &users)
	fmt.Println(users[0])
}

// Private Functions

func registerUser1(db *sql.DB, e *echo.Echo, rec *httptest.ResponseRecorder, h *Handler) error {
	req := httptest.NewRequest(http.MethodPost, "/registration", strings.NewReader(user1JSON))
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

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
