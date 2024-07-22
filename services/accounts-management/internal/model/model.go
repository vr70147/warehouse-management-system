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
	PlanID                uint           `json:"plan_id"` // Foreign key to a subscription plan
	Plan                  Plan           `json:"plan"`    // Associated plan
	IsActive              bool           `json:"is_active" gorm:"default:true"`
	Metadata              string         `json:"metadata" gorm:"type:json"`    // Additional JSON-encoded metadata
	Preferences           string         `json:"preferences" gorm:"type:json"` // JSON-encoded user preferences
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index"`
}

type Plan struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	Phone      string         `json:"phone" gorm:"unique; not null"`
	Street     string         `json:"street"`
	City       string         `json:"city"`
	Password   string         `json:"password" gorm:"not null"`
	IsAdmin    bool           `json:"is_admin" gorm:"default: false"`
	AccountID  uint           `json:"account_id"`
}

type ErrorResponse struct {
	Error string `json:"message"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
