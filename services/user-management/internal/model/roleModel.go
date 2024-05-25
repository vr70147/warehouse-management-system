package model

import (
	"time"
)

type Role struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Role         string    `gorm:"unique:not null"`
	Description  string
	Permission   string     `gorm:"type:jsonb"`
	IsActive     bool       `gorm:"default:true"`
	Users        []User     `gorm:"foreignKey:RoleID"`
	DepartmentID uint       `gorm:"not null"`
	Department   Department `gorm:"foreignKey:DepartmentID"`
}
