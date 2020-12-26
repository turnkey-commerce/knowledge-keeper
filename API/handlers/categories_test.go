package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	echo "github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/turnkey-commerce/knowledge-keeper/models"
)

// See users_test.go for init and TestMain functions
var (
	category1Name       = "New Category"
	category1UpdateName = "Updated Category"
)

// TestCanCreateCategory tests that a new category can be created
func TestCanCreateCategory(t *testing.T) {
	// Setup
	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	rec := httptest.NewRecorder()
	e := echo.New()
	h := NewHandler(db, "secret")

	// First register the user.
	assert.NoError(t, registerUser(db, e, rec, h, user1JSON))

	// Login the user
	rec = httptest.NewRecorder()
	assert.NoError(t, loginUser1(db, e, rec, h, user1LoginJSON))
	token := getLoginToken(rec)

	rec = httptest.NewRecorder()

	if assert.NoError(t, createCategory(db, e, rec, h, token)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		result := make(map[string]interface{})
		json.Unmarshal([]byte(rec.Body.String()), &result)
		resultName := result["name"]

		assert.Equal(t, category1Name, resultName)
	}
}

// TestCanUpdateCategory
func TestCanUpdateCategory(t *testing.T) {
	// Setup
	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	rec := httptest.NewRecorder()
	e := echo.New()
	h := NewHandler(db, "secret")

	// First register the user.
	assert.NoError(t, registerUser(db, e, rec, h, user1JSON))

	// Login the user
	rec = httptest.NewRecorder()
	assert.NoError(t, loginUser1(db, e, rec, h, user1LoginJSON))
	token := getLoginToken(rec)

	rec = httptest.NewRecorder()
	assert.NoError(t, createCategory(db, e, rec, h, token))
	// Get the ID of the new Category
	var category models.Category
	json.Unmarshal([]byte(rec.Body.String()), &category)
	id := category.CategoryID

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/categories/{id}",
		strings.NewReader(getCategory1JSON(category1UpdateName)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))
	c.Set("user", parseToken(token))

	assert.NoError(t, h.UpdateCategory(c))

	var updatedCategory models.Category
	json.Unmarshal([]byte(rec.Body.String()), &updatedCategory)
	assert.Equal(t, category1UpdateName, updatedCategory.Name)
}

// TestCanViewRecentCategories tests that the registered user can view recent categories.
func TestCanViewRecentCategories(t *testing.T) {
	// Setup
	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	rec := httptest.NewRecorder()
	e := echo.New()
	h := NewHandler(db, "secret")

	// First register the user.
	assert.NoError(t, registerUser(db, e, rec, h, user1JSON))

	// Login the user
	rec = httptest.NewRecorder()
	assert.NoError(t, loginUser1(db, e, rec, h, user1LoginJSON))
	token := getLoginToken(rec)

	// Create the category
	assert.NoError(t, createCategory(db, e, rec, h, token))

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	assert.NoError(t, h.GetRecentCategoriesPaginated(c))
	var categories []models.Category
	json.Unmarshal([]byte(rec.Body.String()), &categories)
	assert.Equal(t, category1Name, categories[0].Name)
}

func createCategory(db *sql.DB, e *echo.Echo, rec *httptest.ResponseRecorder, h *Handler, token string) error {
	req := httptest.NewRequest(http.MethodPost, "/categories",
		strings.NewReader(getCategory1JSON(category1Name)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.Set("user", parseToken(token))
	return h.SaveCategory(c)
}

func getCategory1JSON(name string) string {
	return fmt.Sprintf(
		`{"name": "%s", 
    	"description": "Juke" }`, name)
}
