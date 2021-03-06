// Package models contains the types for schema 'public'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"

	nullable "gopkg.in/guregu/null.v4"
)

// Media represents a row from 'public.media'.
type Media struct {
	MediaID     int64           `json:"media_id"`    // media_id
	Type        MediaType       `json:"type"`        // type
	Title       string          `json:"title"`       // title
	Description nullable.String `json:"description"` // description
	URL         string          `json:"url"`         // url
	CreatedBy   int64           `json:"created_by"`  // created_by
	UpdatedBy   nullable.Int    `json:"updated_by"`  // updated_by
	TopicID     int64           `json:"topic_id"`    // topic_id

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Media exists in the database.
func (m *Media) Exists() bool {
	return m._exists
}

// Deleted provides information if the Media has been deleted from the database.
func (m *Media) Deleted() bool {
	return m._deleted
}

// Insert inserts the Media to the database.
func (m *Media) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if m._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO public.media (` +
		`type, title, description, url, created_by, updated_by, topic_id` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`) RETURNING media_id`

	// run query
	XOLog(sqlstr, m.Type, m.Title, m.Description, m.URL, m.CreatedBy, m.UpdatedBy, m.TopicID)
	err = db.QueryRow(sqlstr, m.Type, m.Title, m.Description, m.URL, m.CreatedBy, m.UpdatedBy, m.TopicID).Scan(&m.MediaID)
	if err != nil {
		return err
	}

	// set existence
	m._exists = true

	return nil
}

// Update updates the Media in the database.
func (m *Media) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !m._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if m._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.media SET (` +
		`type, title, description, url, created_by, updated_by, topic_id` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`) WHERE media_id = $8`

	// run query
	XOLog(sqlstr, m.Type, m.Title, m.Description, m.URL, m.CreatedBy, m.UpdatedBy, m.TopicID, m.MediaID)
	_, err = db.Exec(sqlstr, m.Type, m.Title, m.Description, m.URL, m.CreatedBy, m.UpdatedBy, m.TopicID, m.MediaID)
	return err
}

// Save saves the Media to the database.
func (m *Media) Save(db XODB) error {
	if m.Exists() {
		return m.Update(db)
	}

	return m.Insert(db)
}

// Upsert performs an upsert for Media.
//
// NOTE: PostgreSQL 9.5+ only
func (m *Media) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if m._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.media (` +
		`media_id, type, title, description, url, created_by, updated_by, topic_id` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8` +
		`) ON CONFLICT (media_id) DO UPDATE SET (` +
		`media_id, type, title, description, url, created_by, updated_by, topic_id` +
		`) = (` +
		`EXCLUDED.media_id, EXCLUDED.type, EXCLUDED.title, EXCLUDED.description, EXCLUDED.url, EXCLUDED.created_by, EXCLUDED.updated_by, EXCLUDED.topic_id` +
		`)`

	// run query
	XOLog(sqlstr, m.MediaID, m.Type, m.Title, m.Description, m.URL, m.CreatedBy, m.UpdatedBy, m.TopicID)
	_, err = db.Exec(sqlstr, m.MediaID, m.Type, m.Title, m.Description, m.URL, m.CreatedBy, m.UpdatedBy, m.TopicID)
	if err != nil {
		return err
	}

	// set existence
	m._exists = true

	return nil
}

// Delete deletes the Media from the database.
func (m *Media) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !m._exists {
		return nil
	}

	// if deleted, bail
	if m._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.media WHERE media_id = $1`

	// run query
	XOLog(sqlstr, m.MediaID)
	_, err = db.Exec(sqlstr, m.MediaID)
	if err != nil {
		return err
	}

	// set deleted
	m._deleted = true

	return nil
}

// GetRecentPaginatedMedias returns rows from 'public.media',
// that are paginated by the limit and offset inputs.
func GetRecentPaginatedMedias(db XODB, limit int, offset int) ([]*Media, error) {
	const sqlstr = `SELECT ` +
		`media_id, type, title, description, url, created_by, updated_by, topic_id ` +
		`FROM public.media ` +
		`ORDER BY date_created DESC ` +
		`LIMIT $1 OFFSET $2`

	q, err := db.Query(sqlstr, limit, offset)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*Media
	for q.Next() {
		m := Media{}

		// scan
		err = q.Scan(&m.MediaID, &m.Type, &m.Title, &m.Description, &m.URL, &m.CreatedBy, &m.UpdatedBy, &m.TopicID)
		if err != nil {
			return nil, err
		}

		res = append(res, &m)
	}

	return res, nil
}

// MediaByMediaID retrieves a row from 'public.media' as a Media.
//
// Generated from index 'media_pk'.
func MediaByMediaID(db XODB, mediaID int64) (*Media, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`media_id, type, title, description, url, created_by, updated_by, topic_id ` +
		`FROM public.media ` +
		`WHERE media_id = $1`

	// run query
	XOLog(sqlstr, mediaID)
	m := Media{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, mediaID).Scan(&m.MediaID, &m.Type, &m.Title, &m.Description, &m.URL, &m.CreatedBy, &m.UpdatedBy, &m.TopicID)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// MediaByTitle retrieves a row from 'public.media' as a Media.
//
// Generated from index 'media_title_idx'.
func MediaByTitle(db XODB, title string) ([]*Media, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`media_id, type, title, description, url, created_by, updated_by, topic_id ` +
		`FROM public.media ` +
		`WHERE title = $1`

	// run query
	XOLog(sqlstr, title)
	q, err := db.Query(sqlstr, title)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Media{}
	for q.Next() {
		m := Media{
			_exists: true,
		}

		// scan
		err = q.Scan(&m.MediaID, &m.Type, &m.Title, &m.Description, &m.URL, &m.CreatedBy, &m.UpdatedBy, &m.TopicID)
		if err != nil {
			return nil, err
		}

		res = append(res, &m)
	}

	return res, nil
}

// MediaByTopicIDTitle retrieves a row from 'public.media' as a Media.
//
// Generated from index 'media_topics_title_unique_idx'.
func MediaByTopicIDTitle(db XODB, topicID int64, title string) (*Media, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`media_id, type, title, description, url, created_by, updated_by, topic_id ` +
		`FROM public.media ` +
		`WHERE topic_id = $1 AND title = $2`

	// run query
	XOLog(sqlstr, topicID, title)
	m := Media{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, topicID, title).Scan(&m.MediaID, &m.Type, &m.Title, &m.Description, &m.URL, &m.CreatedBy, &m.UpdatedBy, &m.TopicID)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// MediaByType retrieves a row from 'public.media' as a Media.
//
// Generated from index 'media_type_idx'.
func MediaByType(db XODB, typ MediaType) ([]*Media, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`media_id, type, title, description, url, created_by, updated_by, topic_id ` +
		`FROM public.media ` +
		`WHERE type = $1`

	// run query
	XOLog(sqlstr, typ)
	q, err := db.Query(sqlstr, typ)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Media{}
	for q.Next() {
		m := Media{
			_exists: true,
		}

		// scan
		err = q.Scan(&m.MediaID, &m.Type, &m.Title, &m.Description, &m.URL, &m.CreatedBy, &m.UpdatedBy, &m.TopicID)
		if err != nil {
			return nil, err
		}

		res = append(res, &m)
	}

	return res, nil
}
