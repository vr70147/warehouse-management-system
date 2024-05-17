package model

import (
	"inventory-management/internal/kafka"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Quantity    int `gorm:"not null"`
	Location    string
}

type Order struct {
	gorm.Model
	OrderNumber string `gorm:"unique;not null"`
	UserID      uint
	User        kafka.UserEvent
	Items       []OrderItem `gorm:"many2many:order_items"`
	Status      string      `gorm:"not null"`
}

type OrderItem struct {
	gorm.Model
	OrderID         uint
	InventoryItemID uint
	Quantity        int `gorm:"not null"`
}

type Supplier struct {
	gorm.Model
	Name    string `gorm:"not null"`
	Contact string
	Email   string
	Items   []Item `gorm:"many2many:supplier_items"`
}
