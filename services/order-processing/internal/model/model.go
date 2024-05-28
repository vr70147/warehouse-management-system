package model

import (
	"time"
)

type Order struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	ProductID  uint
	Quantity   uint
	CustomerID uint
	Status     string
}

type ErrorResponse struct {
	Error string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Order   Order  `json:"order"`
}

type SuccessResponses struct {
	Message string  `json:"message"`
	Orders  []Order `json:"orders"`
}
