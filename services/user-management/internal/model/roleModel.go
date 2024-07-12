package model

import (
	"time"
)

type Role struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Role         string     `gorm:"unique:not null"`
	Description  string
	Permissions  []Permission `gorm:"many2many:role_permissions"`
	IsActive     bool         `gorm:"default:true"`
	Users        []User       `gorm:"foreignKey:RoleID"`
	DepartmentID uint         `gorm:"not null"`
	Department   Department   `gorm:"foreignKey:DepartmentID"`
	AccountID    uint
}

type Permission struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}
