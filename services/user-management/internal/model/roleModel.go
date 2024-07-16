package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Role         string         `gorm:"unique:not null"`
	Description  string
	IsActive     bool       `gorm:"default:true"`
	Users        []User     `gorm:"foreignKey:RoleID"`
	DepartmentID uint       `gorm:"not null"`
	Department   Department `gorm:"foreignKey:DepartmentID"`
	AccountID    uint
}
