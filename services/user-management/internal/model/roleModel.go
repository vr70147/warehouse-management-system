package model

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	RoleName    string `gorm:"unique:not null"`
	Description string
	Permission  string `gorm:"type:jsonb"`
	IsActive    bool   `gorm:"default:true"`
	Users       []User `gorm:"foreignKey:RoleID"`
}
