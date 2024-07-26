package model

import (
	"database/sql/driver"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Permission string

const (
	PermissionWorker     Permission = "worker"
	PermissionManager    Permission = "manager"
	PermissionSuperAdmin Permission = "super-admin"
)

func (p *Permission) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*p = Permission(v)
	case string:
		*p = Permission(v)
	default:
		return errors.New("invalid permission value")
	}
	return nil
}

func (p Permission) Value() (driver.Value, error) {
	switch p {
	case PermissionWorker, PermissionManager, PermissionSuperAdmin:
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
	RoleID     uint           `json:"role_id"`
	Role       string         `json:"role" gorm:"foreignKey:RoleID"`
	Permission Permission     `json:"permission" gorm:"not null"`
	Phone      string         `json:"phone" gorm:"unique; not null"`
	Street     string         `json:"street"`
	City       string         `json:"city"`
	Password   string         `json:"password" gorm:"not null"`
	IsAdmin    bool           `json:"is_admin" gorm:"default: false"`
	AccountID  uint           `json:"account_id"`
}

type Account struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	Email          string         `json:"email" gorm:"unique;not null"`
	Name           string         `json:"name" gorm:"not null"`
	Password       string         `json:"-"` // Typically hashed, not returned in JSON
	PhoneNumber    string         `json:"phone_number"`
	CompanyName    string         `json:"company_name"`
	Address        string         `json:"address"`
	City           string         `json:"city"`
	State          string         `json:"state"`
	PostalCode     string         `json:"postal_code"`
	Country        string         `json:"country"`
	BillingEmail   string         `json:"billing_email"`
	BillingAddress string         `json:"billing_address"`
	IsActive       bool           `json:"is_active" gorm:"default:true"`
	Metadata       string         `json:"metadata" gorm:"type:json"`    // Additional JSON-encoded metadata
	Preferences    string         `json:"preferences" gorm:"type:json"` // JSON-encoded user preferences
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

type Activity struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
}
