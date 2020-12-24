// Package models contains the types for schema 'public'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
)

// User represents a row from 'public.users'.
type User struct {
	UserID    int64  `json:"user_id"`    // user_id
	Email     string `json:"email"`      // email
	Hash      string `json:"hash"`       // hash
	FirstName string `json:"first_name"` // first_name
	LastName  string `json:"last_name"`  // last_name
	IsAdmin   bool   `json:"is_admin"`   // is_admin
	IsActive  bool   `json:"is_active"`  // is_active

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the User exists in the database.
func (u *User) Exists() bool {
	return u._exists
}

// Deleted provides information if the User has been deleted from the database.
func (u *User) Deleted() bool {
	return u._deleted
}

// Insert inserts the User to the database.
func (u *User) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if u._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO public.users (` +
		`email, hash, first_name, last_name, is_admin, is_active` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6` +
		`) RETURNING user_id`

	// run query
	XOLog(sqlstr, u.Email, u.Hash, u.FirstName, u.LastName, u.IsAdmin, u.IsActive)
	err = db.QueryRow(sqlstr, u.Email, u.Hash, u.FirstName, u.LastName, u.IsAdmin, u.IsActive).Scan(&u.UserID)
	if err != nil {
		return err
	}

	// set existence
	u._exists = true

	return nil
}

// Update updates the User in the database.
func (u *User) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !u._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if u._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.users SET (` +
		`email, hash, first_name, last_name, is_admin, is_active` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6` +
		`) WHERE user_id = $7`

	// run query
	XOLog(sqlstr, u.Email, u.Hash, u.FirstName, u.LastName, u.IsAdmin, u.IsActive, u.UserID)
	_, err = db.Exec(sqlstr, u.Email, u.Hash, u.FirstName, u.LastName, u.IsAdmin, u.IsActive, u.UserID)
	return err
}

// Save saves the User to the database.
func (u *User) Save(db XODB) error {
	if u.Exists() {
		return u.Update(db)
	}

	return u.Insert(db)
}

// Upsert performs an upsert for User.
//
// NOTE: PostgreSQL 9.5+ only
func (u *User) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if u._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.users (` +
		`user_id, email, hash, first_name, last_name, is_admin, is_active` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`) ON CONFLICT (user_id) DO UPDATE SET (` +
		`user_id, email, hash, first_name, last_name, is_admin, is_active` +
		`) = (` +
		`EXCLUDED.user_id, EXCLUDED.email, EXCLUDED.hash, EXCLUDED.first_name, EXCLUDED.last_name, EXCLUDED.is_admin, EXCLUDED.is_active` +
		`)`

	// run query
	XOLog(sqlstr, u.UserID, u.Email, u.Hash, u.FirstName, u.LastName, u.IsAdmin, u.IsActive)
	_, err = db.Exec(sqlstr, u.UserID, u.Email, u.Hash, u.FirstName, u.LastName, u.IsAdmin, u.IsActive)
	if err != nil {
		return err
	}

	// set existence
	u._exists = true

	return nil
}

// Delete deletes the User from the database.
func (u *User) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !u._exists {
		return nil
	}

	// if deleted, bail
	if u._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.users WHERE user_id = $1`

	// run query
	XOLog(sqlstr, u.UserID)
	_, err = db.Exec(sqlstr, u.UserID)
	if err != nil {
		return err
	}

	// set deleted
	u._deleted = true

	return nil
}

// GetRecentPaginatedUsers returns rows from 'public.users',
// that are paginated by the limit and offset inputs.
func GetRecentPaginatedUsers(db XODB, limit int, offset int) ([]*User, error) {
	const sqlstr = `SELECT ` +
		`user_id, email, hash, first_name, last_name, is_admin, is_active ` +
		`FROM public.users ` +
		`ORDER BY date_created DESC ` +
		`LIMIT $1 OFFSET $2`

	q, err := db.Query(sqlstr, limit, offset)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*User
	for q.Next() {
		u := User{}

		// scan
		err = q.Scan(&u.UserID, &u.Email, &u.Hash, &u.FirstName, &u.LastName, &u.IsAdmin, &u.IsActive)
		if err != nil {
			return nil, err
		}

		res = append(res, &u)
	}

	return res, nil
}

// UsersByEmail retrieves a row from 'public.users' as a User.
//
// Generated from index 'users_email_idx'.
func UsersByEmail(db XODB, email string) ([]*User, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`user_id, email, hash, first_name, last_name, is_admin, is_active ` +
		`FROM public.users ` +
		`WHERE email = $1`

	// run query
	XOLog(sqlstr, email)
	q, err := db.Query(sqlstr, email)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*User{}
	for q.Next() {
		u := User{
			_exists: true,
		}

		// scan
		err = q.Scan(&u.UserID, &u.Email, &u.Hash, &u.FirstName, &u.LastName, &u.IsAdmin, &u.IsActive)
		if err != nil {
			return nil, err
		}

		res = append(res, &u)
	}

	return res, nil
}

// UserByUserID retrieves a row from 'public.users' as a User.
//
// Generated from index 'users_pk'.
func UserByUserID(db XODB, userID int64) (*User, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`user_id, email, hash, first_name, last_name, is_admin, is_active ` +
		`FROM public.users ` +
		`WHERE user_id = $1`

	// run query
	XOLog(sqlstr, userID)
	u := User{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, userID).Scan(&u.UserID, &u.Email, &u.Hash, &u.FirstName, &u.LastName, &u.IsAdmin, &u.IsActive)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
