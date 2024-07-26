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
	Description  string         `json:"description"`
	IsActive     bool           `gorm:"default:true"`
	Users        []User         `gorm:"foreignKey:RoleID"`
	DepartmentID uint           `json:"department_id"`
	Department   Department     `gorm:"foreignKey:DepartmentID"`
	AccountID    uint
}
