// Package models contains the types for schema 'public'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"

	nullable "gopkg.in/guregu/null.v4"
)

// Tag represents a row from 'public.tags'.
type Tag struct {
	Name      string       `json:"name"`       // name
	TagID     int64        `json:"tag_id"`     // tag_id
	CreatedBy int64        `json:"created_by"` // created_by
	UpdatedBy nullable.Int `json:"updated_by"` // updated_by

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Tag exists in the database.
func (t *Tag) Exists() bool {
	return t._exists
}

// Deleted provides information if the Tag has been deleted from the database.
func (t *Tag) Deleted() bool {
	return t._deleted
}

// Insert inserts the Tag to the database.
func (t *Tag) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if t._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO public.tags (` +
		`name, created_by, updated_by` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) RETURNING tag_id`

	// run query
	XOLog(sqlstr, t.Name, t.CreatedBy, t.UpdatedBy)
	err = db.QueryRow(sqlstr, t.Name, t.CreatedBy, t.UpdatedBy).Scan(&t.TagID)
	if err != nil {
		return err
	}

	// set existence
	t._exists = true

	return nil
}

// Update updates the Tag in the database.
func (t *Tag) Update(db XODB) error {
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
	const sqlstr = `UPDATE public.tags SET (` +
		`name, created_by, updated_by` +
		`) = ( ` +
		`$1, $2, $3` +
		`) WHERE tag_id = $4`

	// run query
	XOLog(sqlstr, t.Name, t.CreatedBy, t.UpdatedBy, t.TagID)
	_, err = db.Exec(sqlstr, t.Name, t.CreatedBy, t.UpdatedBy, t.TagID)
	return err
}

// Save saves the Tag to the database.
func (t *Tag) Save(db XODB) error {
	if t.Exists() {
		return t.Update(db)
	}

	return t.Insert(db)
}

// Upsert performs an upsert for Tag.
//
// NOTE: PostgreSQL 9.5+ only
func (t *Tag) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if t._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.tags (` +
		`name, tag_id, created_by, updated_by` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) ON CONFLICT (tag_id) DO UPDATE SET (` +
		`name, tag_id, created_by, updated_by` +
		`) = (` +
		`EXCLUDED.name, EXCLUDED.tag_id, EXCLUDED.created_by, EXCLUDED.updated_by` +
		`)`

	// run query
	XOLog(sqlstr, t.Name, t.TagID, t.CreatedBy, t.UpdatedBy)
	_, err = db.Exec(sqlstr, t.Name, t.TagID, t.CreatedBy, t.UpdatedBy)
	if err != nil {
		return err
	}

	// set existence
	t._exists = true

	return nil
}

// Delete deletes the Tag from the database.
func (t *Tag) Delete(db XODB) error {
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
	const sqlstr = `DELETE FROM public.tags WHERE tag_id = $1`

	// run query
	XOLog(sqlstr, t.TagID)
	_, err = db.Exec(sqlstr, t.TagID)
	if err != nil {
		return err
	}

	// set deleted
	t._deleted = true

	return nil
}

// GetRecentPaginatedTags returns rows from 'public.tags',
// that are paginated by the limit and offset inputs.
func GetRecentPaginatedTags(db XODB, limit int, offset int) ([]*Tag, error) {
	const sqlstr = `SELECT ` +
		`name, tag_id, created_by, updated_by ` +
		`FROM public.tags ` +
		`ORDER BY date_created DESC ` +
		`LIMIT $1 OFFSET $2`

	q, err := db.Query(sqlstr, limit, offset)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*Tag
	for q.Next() {
		t := Tag{}

		// scan
		err = q.Scan(&t.Name, &t.TagID, &t.CreatedBy, &t.UpdatedBy)
		if err != nil {
			return nil, err
		}

		res = append(res, &t)
	}

	return res, nil
}

// TagsByName retrieves a row from 'public.tags' as a Tag.
//
// Generated from index 'tags_name_idx'.
func TagsByName(db XODB, name string) ([]*Tag, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`name, tag_id, created_by, updated_by ` +
		`FROM public.tags ` +
		`WHERE name = $1`

	// run query
	XOLog(sqlstr, name)
	q, err := db.Query(sqlstr, name)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Tag{}
	for q.Next() {
		t := Tag{
			_exists: true,
		}

		// scan
		err = q.Scan(&t.Name, &t.TagID, &t.CreatedBy, &t.UpdatedBy)
		if err != nil {
			return nil, err
		}

		res = append(res, &t)
	}

	return res, nil
}

// TagByTagID retrieves a row from 'public.tags' as a Tag.
//
// Generated from index 'tags_pk'.
func TagByTagID(db XODB, tagID int64) (*Tag, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`name, tag_id, created_by, updated_by ` +
		`FROM public.tags ` +
		`WHERE tag_id = $1`

	// run query
	XOLog(sqlstr, tagID)
	t := Tag{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, tagID).Scan(&t.Name, &t.TagID, &t.CreatedBy, &t.UpdatedBy)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
