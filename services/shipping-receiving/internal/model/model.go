package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	ProductID  uint
	Quantity   uint
	CustomerID uint
	Status     string
	AccountID  uint
}

type Shipping struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	OrderID    uint
	ReceiverID uint
	Status     string
	AccountID  uint
}

type ErrorResponse struct {
	Error string `json:"message"`
}

type SuccessResponse struct {
	Message string   `json:"message"`
	Data    Shipping `json:"shipping"`
}

type SuccessResponses struct {
	Message string     `json:"message"`
	Data    []Shipping `json:"shippings"`
}
