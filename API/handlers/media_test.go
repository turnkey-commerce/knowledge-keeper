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
	media1Title       = "New Media"
	media1UpdateTitle = "Updated Media"
)

// TestCanCreateMedia tests that a new media can be created
func TestCanCreateMedia(t *testing.T) {
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
	rec = httptest.NewRecorder()
	assert.NoError(t, createCategory(db, e, rec, h, token))
	// Get the ID of the new Category
	var category models.Category
	json.Unmarshal([]byte(rec.Body.String()), &category)
	categoryID := category.CategoryID

	// Create the Topic
	rec = httptest.NewRecorder()
	assert.NoError(t, createTopic(db, e, rec, h, token, categoryID))
	// Get the ID of the new Topic
	var topic models.Topic
	json.Unmarshal([]byte(rec.Body.String()), &topic)
	topicID := topic.TopicID

	// Create the Media
	rec = httptest.NewRecorder()
	if assert.NoError(t, createMedia(db, e, rec, h, token, topicID)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		result := make(map[string]interface{})
		json.Unmarshal([]byte(rec.Body.String()), &result)
		resultName := result["title"]

		assert.Equal(t, media1Title, resultName)
	}
}

// TestCanUpdateMedia
func TestCanUpdateMedia(t *testing.T) {
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
	rec = httptest.NewRecorder()
	assert.NoError(t, createCategory(db, e, rec, h, token))
	// Get the ID of the new Category
	var category models.Category
	json.Unmarshal([]byte(rec.Body.String()), &category)
	categoryID := category.CategoryID

	// Create the Topic
	rec = httptest.NewRecorder()
	assert.NoError(t, createTopic(db, e, rec, h, token, categoryID))
	// Get the ID of the new Topic
	var topic models.Topic
	json.Unmarshal([]byte(rec.Body.String()), &topic)
	topicID := topic.TopicID

	// Create the Media
	rec = httptest.NewRecorder()
	assert.NoError(t, createMedia(db, e, rec, h, token, topicID))
	// Get the ID of the new Media
	var media models.Media
	json.Unmarshal([]byte(rec.Body.String()), &media)
	id := media.MediaID

	// Update the Media
	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/media/{id}",
		strings.NewReader(getMedia1JSON(media1UpdateTitle, topicID)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))
	c.Set("user", parseToken(token))

	assert.NoError(t, h.UpdateMedia(c))

	var updatedMedia models.Media
	json.Unmarshal([]byte(rec.Body.String()), &updatedMedia)
	assert.Equal(t, media1UpdateTitle, updatedMedia.Title)
}

// TestCanViewRecentMedia tests that the registered user can view recent media.
func TestCanViewRecentMedia(t *testing.T) {
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
	rec = httptest.NewRecorder()
	assert.NoError(t, createCategory(db, e, rec, h, token))
	// Get the ID of the new Category
	var category models.Category
	json.Unmarshal([]byte(rec.Body.String()), &category)
	categoryID := category.CategoryID

	// Create the Topic
	rec = httptest.NewRecorder()
	assert.NoError(t, createTopic(db, e, rec, h, token, categoryID))
	// Get the ID of the new Topic
	var topic models.Topic
	json.Unmarshal([]byte(rec.Body.String()), &topic)
	topicID := topic.TopicID

	// Create the media
	assert.NoError(t, createMedia(db, e, rec, h, token, topicID))

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/media", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	assert.NoError(t, h.GetRecentMediaPaginated(c))
	var media []models.Media
	json.Unmarshal([]byte(rec.Body.String()), &media)
	assert.Equal(t, media1Title, media[0].Title)
}

// TestCanGetMediaByNameOrID tests that a media can be retrieved by its title or MediaID.
func TestCanGetMediaByNameOrID(t *testing.T) {
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
	rec = httptest.NewRecorder()
	assert.NoError(t, createCategory(db, e, rec, h, token))
	// Get the ID of the new Category
	var category models.Category
	json.Unmarshal([]byte(rec.Body.String()), &category)
	categoryID := category.CategoryID

	// Create the Topic
	rec = httptest.NewRecorder()
	assert.NoError(t, createTopic(db, e, rec, h, token, categoryID))
	// Get the ID of the new Topic
	var topic models.Topic
	json.Unmarshal([]byte(rec.Body.String()), &topic)
	topicID := topic.TopicID

	// Create the media
	rec = httptest.NewRecorder()
	assert.NoError(t, createMedia(db, e, rec, h, token, topicID))
	// Get the ID of the new Media
	var media models.Media
	json.Unmarshal([]byte(rec.Body.String()), &media)
	id := media.MediaID

	// Setup to get the media by title.
	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/media/title/{title}", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.SetPath("/media/title/:title")
	c.SetParamNames("title")
	c.SetParamValues(media1Title)
	// Assert
	assert.NoError(t, h.GetMediaByTitle(c))
	var mediaByTitle []models.Media
	json.Unmarshal([]byte(rec.Body.String()), &mediaByTitle)
	assert.Equal(t, media1Title, mediaByTitle[0].Title)

	// Setup to get the media by id.
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/media/{id}", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetPath("/media/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))
	// Assert
	assert.NoError(t, h.GetMediaByMediaID(c))
	var mediaByID models.Media
	json.Unmarshal([]byte(rec.Body.String()), &mediaByID)
	assert.Equal(t, id, mediaByID.MediaID)
	assert.Equal(t, media1Title, mediaByID.Title)
}

// Private functions

func createMedia(db *sql.DB, e *echo.Echo, rec *httptest.ResponseRecorder, h *Handler, token string, topicID int64) error {
	req := httptest.NewRequest(http.MethodPost, "/media",
		strings.NewReader(getMedia1JSON(media1Title, topicID)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.Set("user", parseToken(token))
	return h.SaveMedia(c)
}

func getMedia1JSON(title string, topicID int64) string {
	return fmt.Sprintf(
		`{"title": "%s", 
		"description": "Media Description",
		"type": "link",
		"url": "https://www.google.com",
		"topic_id": %d}`, title, topicID)
}
