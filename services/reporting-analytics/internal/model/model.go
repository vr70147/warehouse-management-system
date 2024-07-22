package model

import (
	"time"

	"gorm.io/gorm"
)

type SalesReport struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	OrderID    uint           `json:"order_id"`
	ProductID  uint           `json:"product_id"`
	Quantity   uint           `json:"quantity"`
	TotalSales float64        `json:"total_price"`
	Timestamp  time.Time      `json:"timestamp"`
	AccountID  uint           `gorm:"index"`
}

type InventoryLevel struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	ProductID    uint           `json:"product_id"`
	CurrentLevel uint           `json:"current_level"`
	LastUpdated  time.Time      `json:"last_updated"`
	AccountID    uint           `gorm:"index"`
	Quantity     int            `json:"quantity"`
}

type ShippingStatus struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	OrderID     uint           `json:"order_id"`
	Status      string         `json:"status"`
	LastUpdated time.Time      `json:"last_updated"`
	AccountID   uint           `gorm:"index"`
}

type UserActivity struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    uint           `json:"user_id"`
	Activity  string         `json:"activity"`
	Timestamp time.Time      `json:"timestamp"`
	AccountID uint           `gorm:"index"`
	Action    string         `json:"action"`
}
