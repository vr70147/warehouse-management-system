package model

import (
	"time"
)

type Product struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	AccountID   uint       `gorm:"index"` // Foreign key to Account
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	CategoryID  uint       `json:"category_id"`
	Category    Category   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category"`
	SupplierID  uint       `json:"supplier_id"`
	Supplier    Supplier   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"supplier"`
	Stocks      []Stock    `json:"stocks" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

type Stock struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	ProductID uint       `json:"product_id"`
	Product   Product    `json:"product"`
	Quantity  uint       `json:"quantity"`
	Location  string     `json:"location"`
	AccountID uint       `gorm:"index"` // Foreign key to Account
}

type Category struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	AccountID   uint       `gorm:"index"` // Foreign key to Account
}

type Supplier struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Email       string     `json:"email"`
	Contact     string     `json:"contact"`
	AccountID   uint       `gorm:"index"` // Foreign key to Account
}

type Order struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	ProductID  uint       `json:"product_id"`
	Quantity   uint       `json:"quantity"`
	CustomerID uint       `json:"customer_id"`
	Status     string     `json:"status"`
	AccountID  uint       `gorm:"index"` // Foreign key to Account
}

type ErrorResponse struct {
	Error string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type StockResponse struct {
	ID          uint   `json:"id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Location    string `json:"location"`
}

type StocksResponse struct {
	Message string          `json:"message"`
	Stocks  []StockResponse `json:"stocks"`
}

type CategoryResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

// CategoriesResponse represents the response for retrieving categories
type CategoriesResponse struct {
	Message    string             `json:"message"`
	Categories []CategoryResponse `json:"categories"`
}

type SupplierResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Contact     string `json:"contact"`
}

// SuppliersResponse represents the response for retrieving supplier items
type SuppliersResponse struct {
	Message   string             `json:"message"`
	Suppliers []SupplierResponse `json:"suppliers"`
}
