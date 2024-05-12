package model

import (
	"gorm.io/gorm"
)

type Roles struct {
	gorm.Model
	RoleName    string `gorm:"unique:not null"`
	Description string
	Permission  string `gorm:"type:jsonb"`
	IsActive    bool   `gorm:"default:true"`
}
