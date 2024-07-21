package model

import (
	"database/sql/driver"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	AccountID   uint           `gorm:"index"` // Foreign key to Account
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	CategoryID  uint           `json:"category_id"`
	Category    Category       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category"`
	SupplierID  uint           `json:"supplier_id"`
	Supplier    Supplier       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"supplier"`
	Stocks      []Stock        `json:"stocks" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

type Stock struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`
	ProductID         uint           `json:"product_id"`
	Product           Product        `json:"product"`
	Quantity          int            `json:"quantity"`
	Location          string         `json:"location"`
	AccountID         uint           `gorm:"index"` // Foreign key to Account
	LowStockThreshold int            `json:"low_stock_threshold"`
}

type Category struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	AccountID   uint           `gorm:"index"`
	ParentID    *uint          `gorm:"index" json:"parent_id,omitempty"`
}

type Supplier struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Email       string         `json:"email"`
	Contact     string         `json:"contact"`
	AccountID   uint           `gorm:"index"` // Foreign key to Account
}

type Order struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	ProductID  uint           `json:"product_id"`
	Quantity   uint           `json:"quantity"`
	CustomerID uint           `json:"customer_id"`
	Status     string         `json:"status"`
	AccountID  uint           `gorm:"index"` // Foreign key to Account
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
	Permission Permission     `json:"permission" gorm:"not null"`
	Phone      string         `json:"phone" gorm:"unique; not null"`
	Street     string         `json:"street"`
	City       string         `json:"city"`
	Password   string         `json:"password" gorm:"not null"`
	IsAdmin    bool           `json:"is_admin" gorm:"default: false"`
	AccountID  uint           `json:"account_id"`
}

type Permission string

const (
	PermissionWorker  Permission = "worker"
	PermissionManager Permission = "manager"
)

func (p *Permission) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*p = Permission(v)
	case string:
		*p = Permission(v)
	default:
		return errors.New("invalid permission value")
	}
	return nil
}

func (p Permission) Value() (driver.Value, error) {
	switch p {
	case PermissionWorker, PermissionManager:
		return string(p), nil
	}
	return nil, errors.New("invalid permission value")
}

func (p Permission) String() string {
	return string(p)
}

type Role struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Role         string         `gorm:"unique:not null"`
	Description  string
	IsActive     bool       `gorm:"default:true"`
	Users        []User     `gorm:"foreignKey:RoleID"`
	DepartmentID uint       `gorm:"not null"`
	Department   Department `gorm:"foreignKey:DepartmentID"`
	AccountID    uint
}

type Department struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `json:"department"`
	Roles     []Role         `json:"roles"`
	IsActive  bool           `gorm:"default:true"`
	AccountID uint           `json:"account_id"`
}

type ErrorResponse struct {
	Error string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type StockResponse struct {
	Message     string `json:"message"`
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
