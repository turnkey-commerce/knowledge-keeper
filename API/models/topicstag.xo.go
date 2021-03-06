// Package models contains the types for schema 'public'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
)

// TopicsTag represents a row from 'public.topics_tags'.
type TopicsTag struct {
	TopicID int64 `json:"topic_id"` // topic_id
	TagID   int64 `json:"tag_id"`   // tag_id

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the TopicsTag exists in the database.
func (tt *TopicsTag) Exists() bool {
	return tt._exists
}

// Deleted provides information if the TopicsTag has been deleted from the database.
func (tt *TopicsTag) Deleted() bool {
	return tt._deleted
}

// Insert inserts the TopicsTag to the database.
func (tt *TopicsTag) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if tt._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO public.topics_tags (` +
		`topic_id, tag_id` +
		`) VALUES (` +
		`$1, $2` +
		`)`

	// run query
	XOLog(sqlstr, tt.TopicID, tt.TagID)
	_, err = db.Exec(sqlstr, tt.TopicID, tt.TagID)
	if err != nil {
		return err
	}

	// set existence
	tt._exists = true

	return nil
}

// Update statements omitted due to lack of fields other than primary key

// Delete deletes the TopicsTag from the database.
func (tt *TopicsTag) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !tt._exists {
		return nil
	}

	// if deleted, bail
	if tt._deleted {
		return nil
	}

	// sql query with composite primary key
	const sqlstr = `DELETE FROM public.topics_tags  WHERE topic_id = $1 AND tag_id = $2`

	// run query
	XOLog(sqlstr, tt.TopicID, tt.TagID)
	_, err = db.Exec(sqlstr, tt.TopicID, tt.TagID)
	if err != nil {
		return err
	}

	// set deleted
	tt._deleted = true

	return nil
}

// GetRecentPaginatedTopicsTags returns rows from 'public.topics_tags',
// that are paginated by the limit and offset inputs.
func GetRecentPaginatedTopicsTags(db XODB, limit int, offset int) ([]*TopicsTag, error) {
	const sqlstr = `SELECT ` +
		`topic_id, tag_id ` +
		`FROM public.topics_tags ` +
		`ORDER BY date_created DESC ` +
		`LIMIT $1 OFFSET $2`

	q, err := db.Query(sqlstr, limit, offset)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*TopicsTag
	for q.Next() {
		tt := TopicsTag{}

		// scan
		err = q.Scan(&tt.TopicID, &tt.TagID)
		if err != nil {
			return nil, err
		}

		res = append(res, &tt)
	}

	return res, nil
}

// TopicsTagByTopicIDTagID retrieves a row from 'public.topics_tags' as a TopicsTag.
//
// Generated from index 'topics_tags_pk'.
func TopicsTagByTopicIDTagID(db XODB, topicID int64, tagID int64) (*TopicsTag, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`topic_id, tag_id ` +
		`FROM public.topics_tags ` +
		`WHERE topic_id = $1 AND tag_id = $2`

	// run query
	XOLog(sqlstr, topicID, tagID)
	tt := TopicsTag{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, topicID, tagID).Scan(&tt.TopicID, &tt.TagID)
	if err != nil {
		return nil, err
	}

	return &tt, nil
}
