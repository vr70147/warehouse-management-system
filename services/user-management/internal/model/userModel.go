package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	PersonalID string `json:"personal_id" gorm:"unique;not null"`
	Name       string `json:"name" gorm:"unique;not null"`
	Email      string `json:"email" gorm:"unique;not null"`
	Age        int    `json:"age" gorm:"not null"`
	BirthDate  string `json:"birthDate" gorm:"not null"`
	RoleID     uint   `json:"role_id" gorm:"not null"`
	Role       string `json:"role" gorm:"foreignKey:RoleID"`
	Phone      string `json:"phone" gorm:"unique; not null"`
	Street     string `json:"street"`
	City       string `json:"city"`
	Password   string `json:"password" gorm:"not null"`
	IsAdmin    bool   `json:"is_admin" gorm:"default: false"`
}

type UserEvent struct {
	EventType string `json:"event_type"`
	User      User   `json:"user"`
}
