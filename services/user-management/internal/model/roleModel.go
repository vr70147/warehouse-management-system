package model

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Role         string `gorm:"unique:not null"`
	Description  string
	Permission   string     `gorm:"type:jsonb"`
	IsActive     bool       `gorm:"default:true"`
	Users        []User     `gorm:"foreignKey:RoleID"`
	DepartmentID uint       `gorm:"not null"`
	Department   Department `gorm:"foreignKey:DepartmentID"`
}
