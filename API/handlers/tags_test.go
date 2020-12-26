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
	Tag1Name       = "New Tag"
	Tag1UpdateName = "Updated Tag"
)

// TestCanCreateTag tests that a new Tag can be created
func TestCanCreateTag(t *testing.T) {
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

	if assert.NoError(t, createTag(db, e, rec, h, token)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		result := make(map[string]interface{})
		json.Unmarshal([]byte(rec.Body.String()), &result)
		resultName := result["name"]

		assert.Equal(t, Tag1Name, resultName)
	}
}

// TestCanUpdateTag
func TestCanUpdateTag(t *testing.T) {
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
	assert.NoError(t, createTag(db, e, rec, h, token))
	// Get the ID of the new Tag
	var Tag models.Tag
	json.Unmarshal([]byte(rec.Body.String()), &Tag)
	id := Tag.TagID

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/tags/{id}",
		strings.NewReader(getTag1JSON(Tag1UpdateName)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))
	c.Set("user", parseToken(token))

	assert.NoError(t, h.UpdateTag(c))

	var updatedTag models.Tag
	json.Unmarshal([]byte(rec.Body.String()), &updatedTag)
	assert.Equal(t, Tag1UpdateName, updatedTag.Name)
}

// TestCanViewRecentTags tests that the registered user can view recent tags.
func TestCanViewRecentTags(t *testing.T) {
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

	// Create the Tag
	assert.NoError(t, createTag(db, e, rec, h, token))

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tags", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	assert.NoError(t, h.GetRecentTagsPaginated(c))
	var tags []models.Tag
	json.Unmarshal([]byte(rec.Body.String()), &tags)
	assert.Equal(t, Tag1Name, tags[0].Name)
}

// TestCanGetTagByNameOrID tests that a Tag can be retrieved by its name or TagID.
func TestCanGetTagByNameOrID(t *testing.T) {
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

	// Create the Tag
	rec = httptest.NewRecorder()
	assert.NoError(t, createTag(db, e, rec, h, token))
	// Get the ID of the new Tag
	var Tag models.Tag
	json.Unmarshal([]byte(rec.Body.String()), &Tag)
	id := Tag.TagID

	// Setup to get the Tag by name.
	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tags/name/{name}", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.SetPath("/tags/name/:name")
	c.SetParamNames("name")
	c.SetParamValues(Tag1Name)
	// Assert
	assert.NoError(t, h.GetTagByName(c))
	var TagByName models.Tag
	json.Unmarshal([]byte(rec.Body.String()), &TagByName)
	assert.Equal(t, Tag1Name, TagByName.Name)

	// Setup to get the Tag by id.
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/tags/{id}", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetPath("/tags/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))
	// Assert
	assert.NoError(t, h.GetTagByTagID(c))
	var TagById models.Tag
	json.Unmarshal([]byte(rec.Body.String()), &TagById)
	assert.Equal(t, id, TagById.TagID)
	assert.Equal(t, Tag1Name, TagById.Name)
}

// Private functions

func createTag(db *sql.DB, e *echo.Echo, rec *httptest.ResponseRecorder, h *Handler, token string) error {
	req := httptest.NewRequest(http.MethodPost, "/tags",
		strings.NewReader(getTag1JSON(Tag1Name)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.Set("user", parseToken(token))
	return h.SaveTag(c)
}

func getTag1JSON(name string) string {
	return fmt.Sprintf(
		`{"name": "%s"}`, name)
}
