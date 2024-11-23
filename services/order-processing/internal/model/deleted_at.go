package model

import (
	"time"

	"gorm.io/gorm"
)

// DeletedAt is a custom type to handle soft delete timestamps.
type DeletedAt struct {
	gorm.DeletedAt
}

// MarshalJSON customizes the JSON representation of DeletedAt.
func (d DeletedAt) MarshalJSON() ([]byte, error) {
	if d.Valid {
		return []byte(`"` + d.Time.Format(time.RFC3339) + `"`), nil
	}
	return []byte("null"), nil
}
