package database_test

import (
	"database/sql"
	"log"
	"testing"

	txdb "github.com/DATA-DOG/go-txdb"
	_ "github.com/lib/pq"
	"github.com/turnkey-commerce/knowledge-keeper/models"
	"github.com/xo/dburl"
)

func init() {
	u, err := dburl.Parse("postgres://knowledge-keeper:knowledge-keeper@localhost/knowledge-keeper?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	txdb.Register("knowledge", "postgres", u.DSN)
}

// TestCreateUserAndCategories tests creating a user and adding
func TestCreateUserAndCategories(t *testing.T) {
	categoryName := "Special"
	userEmail := "test@test.com"

	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	createUserAndCategory(db, userEmail, categoryName)

	usersByEmail, err := models.UsersByEmail(db, userEmail)
	checkError(err)
	if usersByEmail[0].Email != userEmail {
		t.Error("Created user not correct:\n", usersByEmail[0])
	}

	categoriesByName, err := models.CategoriesByName(db, categoryName)
	checkError(err)
	if categoriesByName[0].Name != categoryName {
		t.Error("Created category not correct:\n", categoriesByName[0])
	}
}

func TestCreateTopicAndTags(t *testing.T) {
	categoryName := "Special"
	userEmail := "test@test.com"
	topicTitle := "Test Topic"
	tagName1 := "Tag 1"
	tagName2 := "Tag 2"

	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	createUserAndCategory(db, userEmail, categoryName)

	categoriesByName, err := models.CategoriesByName(db, categoryName)
	checkError(err)
	usersByEmail, err := models.UsersByEmail(db, userEmail)
	checkError(err)
	topic := models.Topic{
		Title:       topicTitle,
		Description: sql.NullString{String: "Topic Description"},
		CategoryID:  categoriesByName[0].CategoryID,
		CreatedBy:   usersByEmail[0].UserID,
	}

	err = topic.Save(db)
	checkError(err)

	tag1 := models.Tag{
		Name:      tagName1,
		CreatedBy: usersByEmail[0].UserID,
	}

	tag2 := models.Tag{
		Name:      tagName2,
		CreatedBy: usersByEmail[0].UserID,
	}

	err = tag1.Save(db)
	checkError(err)

	err = tag2.Save(db)
	checkError(err)

	topicTag1 := models.TopicsTag{
		TopicID: topic.TopicID,
		TagID:   tag1.TagID,
	}

	err = topicTag1.Insert(db)
	checkError(err)

	topicTag2 := models.TopicsTag{
		TopicID: topic.TopicID,
		TagID:   tag2.TagID,
	}

	err = topicTag2.Insert(db)
	checkError(err)

	tagsByTopic, err := models.TopicsTagsByTopicID(db, topic.TopicID)
	checkError(err)

	if !topicsTagsContains(tagsByTopic, tag1.TagID) {
		t.Error("Topics Tag Does not contain Tag 1:\n", tag1)
	}

	if !topicsTagsContains(tagsByTopic, tag2.TagID) {
		t.Error("Topics Tag Does not contain Tag 2:\n", tag2)
	}

	topicsByTitle, err := models.TopicsByTitle(db, topicTitle)
	checkError(err)

	if topicsByTitle[0].Title != topicTitle {
		t.Error("Created topic not correct:\n", topicsByTitle[0])
	}
}

func createUserAndCategory(db *sql.DB, userEmail string, categoryName string) {
	user := models.User{
		Email:     userEmail,
		FirstName: "Jack",
		LastName:  "Test",
	}

	err := user.Save(db)
	checkError(err)

	category := models.Category{
		Name:        categoryName,
		Description: sql.NullString{String: "Special Description"},
		CreatedBy:   user.UserID,
	}

	err = category.Save(db)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func topicsTagsContains(tt []*models.TopicsTag, tagID int64) bool {
	for _, t := range tt {
		if t.TagID == tagID {
			return true
		}
	}
	return false
}
