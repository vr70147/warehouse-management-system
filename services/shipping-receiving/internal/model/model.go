package model

import (
	"time"

	"gorm.io/gorm"
)

type Shipping struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	OrderID      uint           `json:"order_id"`
	ReceiverID   uint           `json:"receiver_id"`
	Status       string         `json:"status"`
	AccountID    uint           `json:"account_id"`
	ShippingDate time.Time      `json:"shipping_date"`
}

type OrderEvent struct {
	OrderID   uint   `json:"order_id"`
	ProductID uint   `json:"product_id"`
	Quantity  uint   `json:"quantity"`
	Action    string `json:"action"` // can be "create", "cancel", "ship"
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
