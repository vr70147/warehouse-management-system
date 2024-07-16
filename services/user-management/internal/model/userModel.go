package model

import (
	"database/sql/driver"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Permission string

const (
	PermissionWorker  Permission = "worker"
	PermissionManager Permission = "manager"
)

func (p *Permission) Scan(value interface{}) error {
	*p = Permission(value.([]byte))
	return nil
}

func (p Permission) Value() (driver.Value, error) {
	switch p {
	case PermissionWorker, PermissionManager:
		return string(p), nil
	}
	return nil, errors.New("invalid permission value")
}

func (p Permission) String() string {
	return string(p)
}

type User struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	PersonalID string         `json:"personal_id" gorm:"unique;not null"`
	Name       string         `json:"name" gorm:"unique;not null"`
	Email      string         `json:"email" gorm:"unique;not null"`
	Age        int            `json:"age" gorm:"not null"`
	BirthDate  string         `json:"birthDate" gorm:"not null"`
	RoleID     uint           `json:"role_id" gorm:"not null"`
	Role       string         `json:"role" gorm:"foreignKey:RoleID"`
	Permission Permission     `json:"permission" gorm:"not null"`
	Phone      string         `json:"phone" gorm:"unique; not null"`
	Street     string         `json:"street"`
	City       string         `json:"city"`
	Password   string         `json:"password" gorm:"not null"`
	IsAdmin    bool           `json:"is_admin" gorm:"default: false"`
	AccountID  uint           `json:"account_id"`
}
