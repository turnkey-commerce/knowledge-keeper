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
	nullable "gopkg.in/guregu/null.v4"
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

	userID, categoryID := createUserAndCategory(db, userEmail, categoryName)
	topicID1, topicID2 := createTopics(db, topicTitle1, topicTitle2, categoryID, userID)

	tag1 := models.Tag{
		Name:      tagName1,
		CreatedBy: userID,
	}

	tag2 := models.Tag{
		Name:      tagName2,
		CreatedBy: userID,
	}

	err = tag1.Save(db)
	checkError(err)

	err = tag2.Save(db)
	checkError(err)

	topic1Tag1 := models.TopicsTag{
		TopicID: topicID1,
		TagID:   tag1.TagID,
	}

	err = topic1Tag1.Insert(db)
	checkError(err)

	topic2Tag1 := models.TopicsTag{
		TopicID: topicID2,
		TagID:   tag1.TagID,
	}

	err = topic2Tag1.Insert(db)
	checkError(err)

	topic1Tag2 := models.TopicsTag{
		TopicID: topicID1,
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

	tagsByTopic, err := queries.TagsByTopicID(db, topicID1)
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

func TestCreateRelatedTopics(t *testing.T) {
	categoryName := "Special"
	userEmail := "test@test.com"
	topicTitle1 := "A Test Topic"
	topicTitle2 := "B Test Topic"

	db, err := sql.Open("knowledge", "identifier")
	defer db.Close()
	checkError(err)

	userID, categoryID := createUserAndCategory(db, userEmail, categoryName)
	topicID1, topicID2 := createTopics(db, topicTitle1, topicTitle2, categoryID, userID)

	relatedTopic := models.RelatedTopic{
		TopicID:        topicID1,
		RelatedTopicID: topicID2,
	}

	err = relatedTopic.Insert(db)
	checkError(err)

	relatedTopics, err := queries.RelatedTopicsByTopicId(db, topicID1)
	if relatedTopics[0].Title.String != topicTitle2 {
		t.Error("Wrong related topic:\n", relatedTopics[0])
	}

}

func createUserAndCategory(db *sql.DB, userEmail string, categoryName string) (int64, int64) {
	user := models.User{
		Email:     userEmail,
		FirstName: "Jack",
		LastName:  "Test",
		IsAdmin:   true,
		Hash:      "testHash",
	}

	err := user.Save(db)
	checkError(err)

	category := models.Category{
		Name:        categoryName,
		Description: nullable.StringFrom("Special Description"),
		CreatedBy:   user.UserID,
	}

	err = category.Save(db)
	checkError(err)

	return user.UserID, category.CategoryID
}

func createTopics(db *sql.DB, topicTitle1 string, topicTitle2 string,
	categoryID int64, userID int64) (int64, int64) {
	topic1 := models.Topic{
		Title:       topicTitle1,
		Description: nullable.StringFrom("Topic1 Description"),
		CategoryID:  categoryID,
		CreatedBy:   userID,
	}

	topic2 := models.Topic{
		Title:       topicTitle2,
		Description: nullable.StringFrom("Topic2 Description"),
		CategoryID:  categoryID,
		CreatedBy:   userID,
	}

	err := topic1.Save(db)
	checkError(err)

	err = topic2.Save(db)
	checkError(err)

	return topic1.TopicID, topic2.TopicID
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
