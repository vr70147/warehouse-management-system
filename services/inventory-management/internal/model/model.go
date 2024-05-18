package model

import (
	"gorm.io/gorm"
)

type InventoryItem struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Quantity    int `gorm:"not null"`
	LocationID  uint
	Location    string
}

type Location struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string `gorm:"not null"`
}

type StockMovement struct {
	gorm.Model
	InventoryItemID uint
	InventoryItem   InventoryItem
	MovementType    string
	Quantity        int
	UserID          uint
	Description     string
}

type Supplier struct {
	gorm.Model
	Name    string `gorm:"not null"`
	Contact string
	Email   string
	Items   []InventoryItem `gorm:"many2many:supplier_items"`
}

type User struct {
	ID      int
	Name    string
	Email   string
	IsAdmin bool
}

type UserEvent struct {
	EventType string `json:"event_type"`
	User      User   `json:"user"`
}
