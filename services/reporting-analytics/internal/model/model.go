package model

import (
	"time"
)

type SalesReport struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	OrderID    uint       `json:"order_id"`
	ProductID  uint       `json:"product_id"`
	Quantity   uint       `json:"quantity"`
	TotalPrice float64    `json:"total_price"`
	Timestamp  time.Time  `json:"timestamp"`
}

type InventoryLevel struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	ProductID    uint       `json:"product_id"`
	CurrentLevel uint       `json:"current_level"`
	LastUpdated  time.Time  `json:"last_updated"`
}

type ShippingStatus struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	OrderID     uint       `json:"order_id"`
	Status      string     `json:"status"`
	LastUpdated time.Time  `json:"last_updated"`
}

type UserActivity struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	UserID    uint       `json:"user_id"`
	Activity  string     `json:"activity"`
	Timestamp time.Time  `json:"timestamp"`
}
