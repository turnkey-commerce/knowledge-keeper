// Package queries contains additional queries not generated by gendal.
package queries

import (
	"github.com/turnkey-commerce/knowledge-keeper/models"
)

// TopicsByTagID retrieves all Topics with a given TagID.
//
func TopicsByTagID(db models.XODB, tagID int64) ([]*models.TagTopicsView, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`topic_id, tag_id, category_id, title, description, created_by, updated_by ` +
		`FROM public.tag_topics_view ` +
		`WHERE tag_id = $1` +
		`ORDER BY title`

	// run query
	models.XOLog(sqlstr, tagID)
	q, err := db.Query(sqlstr, tagID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*models.TagTopicsView{}
	for q.Next() {
		t := models.TagTopicsView{}

		// scan
		err = q.Scan(&t.TopicID, &t.TagID, &t.CategoryID, &t.Title, &t.Description, &t.CreatedBy, &t.UpdatedBy)
		if err != nil {
			return nil, err
		}

		res = append(res, &t)
	}

	return res, nil
}

// TagsByTopicID retrieves all Tags with a given TopicId.
//
func TagsByTopicID(db models.XODB, topicID int64) ([]*models.TopicsTagsView, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`name, topic_id, tag_id, created_by, updated_by ` +
		`FROM public.topics_tags_view ` +
		`WHERE topic_id = $1` +
		`ORDER BY name`

	// run query
	models.XOLog(sqlstr, topicID)
	q, err := db.Query(sqlstr, topicID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*models.TopicsTagsView{}
	for q.Next() {
		t := models.TopicsTagsView{}

		// scan
		err = q.Scan(&t.Name, &t.TopicID, &t.TagID, &t.CreatedBy, &t.UpdatedBy)
		if err != nil {
			return nil, err
		}

		res = append(res, &t)
	}

	return res, nil
}

// RelatedTopicsByTopicId retrieves all Topics related to a given TopicId.
//
func RelatedTopicsByTopicId(db models.XODB, tagID int64) ([]*models.RelatedTopicsView, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`topic_id, related_topic_id, category_id, title, description, created_by, updated_by ` +
		`FROM public.related_topics_view ` +
		`WHERE topic_id = $1` +
		`ORDER BY title`

	// run query
	models.XOLog(sqlstr, tagID)
	q, err := db.Query(sqlstr, tagID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*models.RelatedTopicsView{}
	for q.Next() {
		t := models.RelatedTopicsView{}

		// scan
		err = q.Scan(&t.TopicID, &t.RelatedTopicID, &t.CategoryID, &t.Title, &t.Description, &t.CreatedBy, &t.UpdatedBy)
		if err != nil {
			return nil, err
		}

		res = append(res, &t)
	}

	return res, nil
}