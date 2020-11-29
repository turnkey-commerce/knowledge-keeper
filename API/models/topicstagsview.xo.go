// Package models contains the types for schema 'public'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
)

// TopicsTagsView represents a row from 'public.topics_tags_view'.
type TopicsTagsView struct {
	TopicID   sql.NullInt64  `json:"topic_id"`   // topic_id
	Name      sql.NullString `json:"name"`       // name
	TagID     sql.NullInt64  `json:"tag_id"`     // tag_id
	CreatedBy sql.NullInt64  `json:"created_by"` // created_by
	UpdatedBy sql.NullInt64  `json:"updated_by"` // updated_by
}