package model

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID                    uint           `gorm:"primarykey" json:"id"`
	Email                 string         `json:"email" gorm:"unique;not null"`
	Name                  string         `json:"name" gorm:"not null"`
	Password              string         `json:"-"`
	PhoneNumber           string         `json:"phone_number"`
	CompanyName           string         `json:"company_name"`
	Address               string         `json:"address"`
	City                  string         `json:"city"`
	State                 string         `json:"state"`
	PostalCode            string         `json:"postal_code"`
	Country               string         `json:"country"`
	BillingEmail          string         `json:"billing_email"`
	BillingAddress        string         `json:"billing_address"`
	SubscriptionExpiresAt time.Time      `json:"subscription_expires_at"`
	IsActive              bool           `json:"is_active" gorm:"default:true"`
	Metadata              string         `json:"metadata" gorm:"type:json"`    // Additional JSON-encoded metadata
	Preferences           string         `json:"preferences" gorm:"type:json"` // JSON-encoded user preferences
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index"`
}

type User struct {
	ID         uint   `json:"id"`
	PersonalID string `json:"personal_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	RoleID     uint   `json:"role_id"`
	Role       string `json:"role"`
	Permission string `json:"permission"`
	IsAdmin    bool   `json:"is_admin"`
	AccountID  uint   `json:"account_id"`
}

type ErrorResponse struct {
	Error string `json:"message"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
