package database_test

import (
	"database/sql"
	"log"
	"testing"

	txdb "github.com/DATA-DOG/go-txdb"
	_ "github.com/lib/pq"
	"github.com/turnkey-commerce/knowledge-keeper/models"
	"github.com/turnkey-commerce/knowledge-keeper/models/queries"
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
	topicTitle1 := "A Test Topic"
	topicTitle2 := "B Test Topic"
	tagName1 := "A Tag"
	tagName2 := "B Tag"

	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	createUserAndCategory(db, userEmail, categoryName)

	categoriesByName, err := models.CategoriesByName(db, categoryName)
	checkError(err)
	usersByEmail, err := models.UsersByEmail(db, userEmail)
	checkError(err)
	topic1 := models.Topic{
		Title:       topicTitle1,
		Description: sql.NullString{String: "Topic1 Description"},
		CategoryID:  categoriesByName[0].CategoryID,
		CreatedBy:   usersByEmail[0].UserID,
	}

	topic2 := models.Topic{
		Title:       topicTitle2,
		Description: sql.NullString{String: "Topic2 Description"},
		CategoryID:  categoriesByName[0].CategoryID,
		CreatedBy:   usersByEmail[0].UserID,
	}

	err = topic1.Save(db)
	checkError(err)

	err = topic2.Save(db)
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

	topic1Tag1 := models.TopicsTag{
		TopicID: topic1.TopicID,
		TagID:   tag1.TagID,
	}

	err = topic1Tag1.Insert(db)
	checkError(err)

	topic2Tag1 := models.TopicsTag{
		TopicID: topic2.TopicID,
		TagID:   tag1.TagID,
	}

	err = topic2Tag1.Insert(db)
	checkError(err)

	topic1Tag2 := models.TopicsTag{
		TopicID: topic1.TopicID,
		TagID:   tag2.TagID,
	}

	err = topic1Tag2.Insert(db)
	checkError(err)

	topicsByTag, err := queries.TopicsByTagID(db, tag1.TagID)
	checkError(err)

	if topicsByTag[0].Title.String != topicTitle1 {
		t.Error("TopicsByTag wrong topic:\n", topicsByTag[0])
	}

	if topicsByTag[1].Title.String != topicTitle2 {
		t.Error("TopicsByTag wrong topic:\n", topicsByTag[1])
	}

	tagsByTopic, err := queries.TagsByTopicID(db, topic1.TopicID)
	checkError(err)

	if tagsByTopic[0].Name.String != tagName1 {
		t.Error("TagsByTopic wrong tag:\n", tagsByTopic[0])
	}

	if tagsByTopic[1].Name.String != tagName2 {
		t.Error("TagsByTopic wrong topic:\n", tagsByTopic[1])
	}

	topicsByTitle, err := models.TopicsByTitle(db, topicTitle1)
	checkError(err)

	if topicsByTitle[0].Title != topicTitle1 {
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
