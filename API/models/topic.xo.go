// Package models contains the types for schema 'public'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

// Topic represents a row from 'public.topics'.
type Topic struct {
	TopicID     int64          `json:"topic_id"`     // topic_id
	CategoryID  int64          `json:"category_id"`  // category_id
	Title       string         `json:"title"`        // title
	Description sql.NullString `json:"description"`  // description
	CreatedBy   int64          `json:"created_by"`   // created_by
	UpdatedBy   sql.NullInt64  `json:"updated_by"`   // updated_by
	DateCreated time.Time      `json:"date_created"` // date_created
	DateUpdated pq.NullTime    `json:"date_updated"` // date_updated

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Topic exists in the database.
func (t *Topic) Exists() bool {
	return t._exists
}

// Deleted provides information if the Topic has been deleted from the database.
func (t *Topic) Deleted() bool {
	return t._deleted
}

// Insert inserts the Topic to the database.
func (t *Topic) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if t._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO public.topics (` +
		`category_id, title, description, created_by, updated_by, date_created, date_updated` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`) RETURNING topic_id`

	// run query
	XOLog(sqlstr, t.CategoryID, t.Title, t.Description, t.CreatedBy, t.UpdatedBy, t.DateCreated, t.DateUpdated)
	err = db.QueryRow(sqlstr, t.CategoryID, t.Title, t.Description, t.CreatedBy, t.UpdatedBy, t.DateCreated, t.DateUpdated).Scan(&t.TopicID)
	if err != nil {
		return err
	}

	// set existence
	t._exists = true

	return nil
}

// Update updates the Topic in the database.
func (t *Topic) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !t._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if t._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.topics SET (` +
		`category_id, title, description, created_by, updated_by, date_created, date_updated` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`) WHERE topic_id = $8`

	// run query
	XOLog(sqlstr, t.CategoryID, t.Title, t.Description, t.CreatedBy, t.UpdatedBy, t.DateCreated, t.DateUpdated, t.TopicID)
	_, err = db.Exec(sqlstr, t.CategoryID, t.Title, t.Description, t.CreatedBy, t.UpdatedBy, t.DateCreated, t.DateUpdated, t.TopicID)
	return err
}

// Save saves the Topic to the database.
func (t *Topic) Save(db XODB) error {
	if t.Exists() {
		return t.Update(db)
	}

	return t.Insert(db)
}

// Upsert performs an upsert for Topic.
//
// NOTE: PostgreSQL 9.5+ only
func (t *Topic) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if t._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.topics (` +
		`topic_id, category_id, title, description, created_by, updated_by, date_created, date_updated` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8` +
		`) ON CONFLICT (topic_id) DO UPDATE SET (` +
		`topic_id, category_id, title, description, created_by, updated_by, date_created, date_updated` +
		`) = (` +
		`EXCLUDED.topic_id, EXCLUDED.category_id, EXCLUDED.title, EXCLUDED.description, EXCLUDED.created_by, EXCLUDED.updated_by, EXCLUDED.date_created, EXCLUDED.date_updated` +
		`)`

	// run query
	XOLog(sqlstr, t.TopicID, t.CategoryID, t.Title, t.Description, t.CreatedBy, t.UpdatedBy, t.DateCreated, t.DateUpdated)
	_, err = db.Exec(sqlstr, t.TopicID, t.CategoryID, t.Title, t.Description, t.CreatedBy, t.UpdatedBy, t.DateCreated, t.DateUpdated)
	if err != nil {
		return err
	}

	// set existence
	t._exists = true

	return nil
}

// Delete deletes the Topic from the database.
func (t *Topic) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !t._exists {
		return nil
	}

	// if deleted, bail
	if t._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.topics WHERE topic_id = $1`

	// run query
	XOLog(sqlstr, t.TopicID)
	_, err = db.Exec(sqlstr, t.TopicID)
	if err != nil {
		return err
	}

	// set deleted
	t._deleted = true

	return nil
}

// Category returns the Category associated with the Topic's CategoryID (category_id).
//
// Generated from foreign key 'topics_category_fk'.
func (t *Topic) Category(db XODB) (*Category, error) {
	return CategoryByCategoryID(db, t.CategoryID)
}

// UserByCreatedBy returns the User associated with the Topic's CreatedBy (created_by).
//
// Generated from foreign key 'topics_created_by_fk'.
func (t *Topic) UserByCreatedBy(db XODB) (*User, error) {
	return UserByUserID(db, t.CreatedBy)
}

// UserByUpdatedBy returns the User associated with the Topic's UpdatedBy (updated_by).
//
// Generated from foreign key 'topics_updated_by_fk'.
func (t *Topic) UserByUpdatedBy(db XODB) (*User, error) {
	return UserByUserID(db, t.UpdatedBy.Int64)
}

// TopicByTopicID retrieves a row from 'public.topics' as a Topic.
//
// Generated from index 'topics_pk'.
func TopicByTopicID(db XODB, topicID int64) (*Topic, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`topic_id, category_id, title, description, created_by, updated_by, date_created, date_updated ` +
		`FROM public.topics ` +
		`WHERE topic_id = $1`

	// run query
	XOLog(sqlstr, topicID)
	t := Topic{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, topicID).Scan(&t.TopicID, &t.CategoryID, &t.Title, &t.Description, &t.CreatedBy, &t.UpdatedBy, &t.DateCreated, &t.DateUpdated)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// TopicsByTitle retrieves a row from 'public.topics' as a Topic.
//
// Generated from index 'topics_title_idx'.
func TopicsByTitle(db XODB, title string) ([]*Topic, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`topic_id, category_id, title, description, created_by, updated_by, date_created, date_updated ` +
		`FROM public.topics ` +
		`WHERE title = $1`

	// run query
	XOLog(sqlstr, title)
	q, err := db.Query(sqlstr, title)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Topic{}
	for q.Next() {
		t := Topic{
			_exists: true,
		}

		// scan
		err = q.Scan(&t.TopicID, &t.CategoryID, &t.Title, &t.Description, &t.CreatedBy, &t.UpdatedBy, &t.DateCreated, &t.DateUpdated)
		if err != nil {
			return nil, err
		}

		res = append(res, &t)
	}

	return res, nil
}
