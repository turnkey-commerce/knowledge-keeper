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
	note1Title       = "New Note"
	note1UpdateTitle = "Updated Note"
)

// TestCanCreateNote tests that a new note can be created
func TestCanCreateNote(t *testing.T) {
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

	// Create the Note
	rec = httptest.NewRecorder()
	if assert.NoError(t, createNote(db, e, rec, h, token, topicID)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		result := make(map[string]interface{})
		json.Unmarshal([]byte(rec.Body.String()), &result)
		resultName := result["title"]

		assert.Equal(t, note1Title, resultName)
	}
}

// TestCanUpdateNote
func TestCanUpdateNote(t *testing.T) {
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

	// Create the Note
	rec = httptest.NewRecorder()
	assert.NoError(t, createNote(db, e, rec, h, token, topicID))
	// Get the ID of the new Note
	var note models.Note
	json.Unmarshal([]byte(rec.Body.String()), &note)
	id := note.NoteID

	// Update the Note
	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/notes/{id}",
		strings.NewReader(getNote1JSON(note1UpdateTitle, topicID)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))
	c.Set("user", parseToken(token))

	assert.NoError(t, h.UpdateNote(c))

	var updatedNote models.Note
	json.Unmarshal([]byte(rec.Body.String()), &updatedNote)
	assert.Equal(t, note1UpdateTitle, updatedNote.Title)
}

// TestCanViewRecentNotes tests that the registered user can view recent notes.
func TestCanViewRecentNotes(t *testing.T) {
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

	// Create the note
	assert.NoError(t, createNote(db, e, rec, h, token, topicID))

	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	assert.NoError(t, h.GetRecentNotesPaginated(c))
	var notes []models.Note
	json.Unmarshal([]byte(rec.Body.String()), &notes)
	assert.Equal(t, note1Title, notes[0].Title)
}

// TestCanGetNoteByNameOrID tests that a note can be retrieved by its title or NoteID.
func TestCanGetNoteByNameOrID(t *testing.T) {
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

	// Create the note
	rec = httptest.NewRecorder()
	assert.NoError(t, createNote(db, e, rec, h, token, topicID))
	// Get the ID of the new Note
	var note models.Note
	json.Unmarshal([]byte(rec.Body.String()), &note)
	id := note.NoteID

	// Setup to get the note by title.
	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notes/title/{title}", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.SetPath("/notes/title/:title")
	c.SetParamNames("title")
	c.SetParamValues(note1Title)
	// Assert
	assert.NoError(t, h.GetNoteByTitle(c))
	var notesByTitle []models.Note
	json.Unmarshal([]byte(rec.Body.String()), &notesByTitle)
	assert.Equal(t, note1Title, notesByTitle[0].Title)

	// Setup to get the note by id.
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/notes/{id}", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c = e.NewContext(req, rec)
	c.SetPath("/notes/:id")
	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", id))
	// Assert
	assert.NoError(t, h.GetNoteByNoteID(c))
	var noteByID models.Note
	json.Unmarshal([]byte(rec.Body.String()), &noteByID)
	assert.Equal(t, id, noteByID.NoteID)
	assert.Equal(t, note1Title, noteByID.Title)
}

// Private functions

func createNote(db *sql.DB, e *echo.Echo, rec *httptest.ResponseRecorder, h *Handler, token string, topicID int64) error {
	req := httptest.NewRequest(http.MethodPost, "/notes",
		strings.NewReader(getNote1JSON(note1Title, topicID)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.Set("user", parseToken(token))
	return h.SaveNote(c)
}

func getNote1JSON(title string, topicID int64) string {
	return fmt.Sprintf(
		`{"title": "%s", 
		"description": "Note Description",
		"topic_id": %d}`, title, topicID)
}
