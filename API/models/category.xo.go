// Package models contains the types for schema 'public'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"

	nullable "gopkg.in/guregu/null.v4"
)

// Category represents a row from 'public.categories'.
type Category struct {
	CategoryID  int64           `json:"category_id"` // category_id
	Name        string          `json:"name"`        // name
	Description nullable.String `json:"description"` // description
	CreatedBy   int64           `json:"created_by"`  // created_by
	UpdatedBy   nullable.Int    `json:"updated_by"`  // updated_by

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Category exists in the database.
func (c *Category) Exists() bool {
	return c._exists
}

// Deleted provides information if the Category has been deleted from the database.
func (c *Category) Deleted() bool {
	return c._deleted
}

// Insert inserts the Category to the database.
func (c *Category) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if c._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO public.categories (` +
		`name, description, created_by, updated_by` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) RETURNING category_id`

	// run query
	XOLog(sqlstr, c.Name, c.Description, c.CreatedBy, c.UpdatedBy)
	err = db.QueryRow(sqlstr, c.Name, c.Description, c.CreatedBy, c.UpdatedBy).Scan(&c.CategoryID)
	if err != nil {
		return err
	}

	// set existence
	c._exists = true

	return nil
}

// Update updates the Category in the database.
func (c *Category) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !c._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if c._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.categories SET (` +
		`name, description, created_by, updated_by` +
		`) = ( ` +
		`$1, $2, $3, $4` +
		`) WHERE category_id = $5`

	// run query
	XOLog(sqlstr, c.Name, c.Description, c.CreatedBy, c.UpdatedBy, c.CategoryID)
	_, err = db.Exec(sqlstr, c.Name, c.Description, c.CreatedBy, c.UpdatedBy, c.CategoryID)
	return err
}

// Save saves the Category to the database.
func (c *Category) Save(db XODB) error {
	if c.Exists() {
		return c.Update(db)
	}

	return c.Insert(db)
}

// Upsert performs an upsert for Category.
//
// NOTE: PostgreSQL 9.5+ only
func (c *Category) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if c._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.categories (` +
		`category_id, name, description, created_by, updated_by` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) ON CONFLICT (category_id) DO UPDATE SET (` +
		`category_id, name, description, created_by, updated_by` +
		`) = (` +
		`EXCLUDED.category_id, EXCLUDED.name, EXCLUDED.description, EXCLUDED.created_by, EXCLUDED.updated_by` +
		`)`

	// run query
	XOLog(sqlstr, c.CategoryID, c.Name, c.Description, c.CreatedBy, c.UpdatedBy)
	_, err = db.Exec(sqlstr, c.CategoryID, c.Name, c.Description, c.CreatedBy, c.UpdatedBy)
	if err != nil {
		return err
	}

	// set existence
	c._exists = true

	return nil
}

// Delete deletes the Category from the database.
func (c *Category) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !c._exists {
		return nil
	}

	// if deleted, bail
	if c._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.categories WHERE category_id = $1`

	// run query
	XOLog(sqlstr, c.CategoryID)
	_, err = db.Exec(sqlstr, c.CategoryID)
	if err != nil {
		return err
	}

	// set deleted
	c._deleted = true

	return nil
}

// GetRecentPaginatedCategorys returns rows from 'public.categories',
// that are paginated by the limit and offset inputs.
func GetRecentPaginatedCategorys(db XODB, limit int, offset int) ([]*Category, error) {
	const sqlstr = `SELECT ` +
		`category_id, name, description, created_by, updated_by ` +
		`FROM public.categories ` +
		`ORDER BY date_created DESC ` +
		`LIMIT $1 OFFSET $2`

	q, err := db.Query(sqlstr, limit, offset)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*Category
	for q.Next() {
		c := Category{}

		// scan
		err = q.Scan(&c.CategoryID, &c.Name, &c.Description, &c.CreatedBy, &c.UpdatedBy)
		if err != nil {
			return nil, err
		}

		res = append(res, &c)
	}

	return res, nil
}

// CategoryByName retrieves a row from 'public.categories' as a Category.
//
// Generated from index 'categories_name_unique_idx'.
func CategoryByName(db XODB, name string) (*Category, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`category_id, name, description, created_by, updated_by ` +
		`FROM public.categories ` +
		`WHERE name = $1`

	// run query
	XOLog(sqlstr, name)
	c := Category{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, name).Scan(&c.CategoryID, &c.Name, &c.Description, &c.CreatedBy, &c.UpdatedBy)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// CategoryByCategoryID retrieves a row from 'public.categories' as a Category.
//
// Generated from index 'categories_pk'.
func CategoryByCategoryID(db XODB, categoryID int64) (*Category, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`category_id, name, description, created_by, updated_by ` +
		`FROM public.categories ` +
		`WHERE category_id = $1`

	// run query
	XOLog(sqlstr, categoryID)
	c := Category{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, categoryID).Scan(&c.CategoryID, &c.Name, &c.Description, &c.CreatedBy, &c.UpdatedBy)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
