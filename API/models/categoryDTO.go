package models

import (
	nullable "gopkg.in/guregu/null.v4"
)

// CategoryDTO provides a DTO for the category and to provide a proper swaggertype.
type CategoryDTO struct {
	Name        string          `json:"name"`                             // name
	Description nullable.String `json:"description" swaggertype:"string"` // description
}
