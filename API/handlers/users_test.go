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
	"github.com/xo/dburl"
)

var (
	email    = "juke@email.com"
	userJSON = fmt.Sprintf(`{"email": "%s", 
    "first_name": "Juke",
    "last_name": "Smith",
    "password": "Test-Password",
    "is_admin": true}`, email)
)

func init() {
	u, err := dburl.Parse("postgres://knowledge-keeper:knowledge-keeper@localhost/knowledge-keeper?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	txdb.Register("knowledge", "postgres", u.DSN)
}

// TestCreateUser tests the route to create a new user.
func TestCreateUser(t *testing.T) {
	// Setup
	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewHandler(db, "secret")

	// Assertions
	if assert.NoError(t, h.SaveUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		result := make(map[string]interface{})
		json.Unmarshal([]byte(rec.Body.String()), &result)
		resultEmail := result["email"]

		assert.Equal(t, email, resultEmail)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
