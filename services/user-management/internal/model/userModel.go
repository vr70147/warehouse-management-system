package model

import (
	"time"
)

type User struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	PersonalID string     `json:"personal_id" gorm:"unique;not null"`
	Name       string     `json:"name" gorm:"unique;not null"`
	Email      string     `json:"email" gorm:"unique;not null"`
	Age        int        `json:"age" gorm:"not null"`
	BirthDate  string     `json:"birthDate" gorm:"not null"`
	RoleID     uint       `json:"role_id" gorm:"not null"`
	Role       string     `json:"role" gorm:"foreignKey:RoleID"`
	Phone      string     `json:"phone" gorm:"unique; not null"`
	Street     string     `json:"street"`
	City       string     `json:"city"`
	Password   string     `json:"password" gorm:"not null"`
	IsAdmin    bool       `json:"is_admin" gorm:"default: false"`
	AccountID  uint       `json:"account_id"`
}
