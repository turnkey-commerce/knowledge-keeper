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
	topic1Title       = "New Topic"
	topic1UpdateTitle = "Updated Topic"
)

// TestCanCreateTopic tests that a new topic can be created
func TestCanCreateTopic(t *testing.T) {
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

	rec = httptest.NewRecorder()

	if assert.NoError(t, createTopic(db, e, rec, h, token, categoryID)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		result := make(map[string]interface{})
		json.Unmarshal([]byte(rec.Body.String()), &result)
		resultName := result["title"]

		assert.Equal(t, topic1Title, resultName)
	}
}

// TestCanUpdateTopic
func TestCanUpdateTopic(t *testing.T) {
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

	rec = httptest.NewRecorder()
	assert.NoError(t, createTopic(db, e, rec, h, token, categoryID))
	// Get the ID of the new Topic
	var topic models.Topic
	json.Unmarshal([]byte(rec.Body.String()), &topic)
	id := topic.TopicID

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/topics/{id}",
		strings.NewReader(getTopic1JSON(topic1UpdateTitle, categoryID)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))
	c.Set("user", parseToken(token))

	assert.NoError(t, h.UpdateTopic(c))

	var updatedTopic models.Topic
	json.Unmarshal([]byte(rec.Body.String()), &updatedTopic)
	assert.Equal(t, topic1UpdateTitle, updatedTopic.Title)
}

// TestCanViewRecentTopics tests that the registered user can view recent topics.
func TestCanViewRecentTopics(t *testing.T) {
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

	// Create the topic
	assert.NoError(t, createTopic(db, e, rec, h, token, categoryID))

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/topics", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	assert.NoError(t, h.GetRecentTopicsPaginated(c))
	var topics []models.Topic
	json.Unmarshal([]byte(rec.Body.String()), &topics)
	assert.Equal(t, topic1Title, topics[0].Title)
}

// TestCanGetTopicByNameOrID tests that a topic can be retrieved by its title or TopicID.
func TestCanGetTopicByNameOrID(t *testing.T) {
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

	// Create the topic
	rec = httptest.NewRecorder()
	assert.NoError(t, createTopic(db, e, rec, h, token, categoryID))
	// Get the ID of the new Topic
	var topic models.Topic
	json.Unmarshal([]byte(rec.Body.String()), &topic)
	id := topic.TopicID

	// Setup to get the topic by title.
	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/topics/title/{title}", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.SetPath("/topics/title/:title")
	c.SetParamNames("title")
	c.SetParamValues(topic1Title)
	// Assert
	assert.NoError(t, h.GetTopicByTitle(c))
	var topicsByTitle []models.Topic
	json.Unmarshal([]byte(rec.Body.String()), &topicsByTitle)
	assert.Equal(t, topic1Title, topicsByTitle[0].Title)

	// Setup to get the topic by id.
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/topics/{id}", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetPath("/topics/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))
	// Assert
	assert.NoError(t, h.GetTopicByTopicID(c))
	var topicById models.Topic
	json.Unmarshal([]byte(rec.Body.String()), &topicById)
	assert.Equal(t, id, topicById.TopicID)
	assert.Equal(t, topic1Title, topicById.Title)
}

// Private functions

func createTopic(db *sql.DB, e *echo.Echo, rec *httptest.ResponseRecorder, h *Handler, token string, categoryID int64) error {
	req := httptest.NewRequest(http.MethodPost, "/topics",
		strings.NewReader(getTopic1JSON(topic1Title, categoryID)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.Set("user", parseToken(token))
	return h.SaveTopic(c)
}

func getTopic1JSON(title string, categoryID int64) string {
	return fmt.Sprintf(
		`{"title": "%s", 
		"description": "Juke",
		"category_id": %d}`, title, categoryID)
}
