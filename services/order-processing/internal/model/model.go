package model

import (
	"database/sql/driver"
	"errors"
	"time"
)

// Order represents an order in the system.
type Order struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" swaggertype:"string" example:"2023-01-01T00:00:00Z"`
	AccountID    uint       `gorm:"index" json:"account_id"`
	ProductID    uint       `json:"product_id"`
	Quantity     uint       `json:"quantity"`
	CustomerID   uint       `json:"customer_id"`
	Status       string     `json:"status"`
	Version      int        `json:"version"`
	ShippingDate time.Time  `json:"shipping_date"`
}

// OrderStatusUpdateRequest represents the payload to update the status of an order.
type OrderStatusUpdateRequest struct {
	Status string `json:"status" binding:"required"`
}

// OrderEvent represents an order event for Kafka.
type OrderEvent struct {
	OrderID   uint   `json:"order_id"`
	ProductID uint   `json:"product_id"`
	Quantity  uint   `json:"quantity"`
	Action    string `json:"action"`
}

// InventoryStatusEvent represents an inventory status event.
type InventoryStatusEvent struct {
	OrderID uint   `json:"order_id"`
	Status  string `json:"status"`
}

// ErrorResponse represents a generic error response.
type ErrorResponse struct {
	Error string `json:"message"`
}

// SuccessResponse represents a generic success response.
type SuccessResponse struct {
	Message string `json:"message"`
	Order   Order  `json:"order"`
}

// SuccessCustomerResponse represents a success response with an order and customer details.
type SuccessCustomerResponse struct {
	Message  string   `json:"message"`
	Order    Order    `json:"order"`
	Customer Customer `json:"customer"`
}

// SuccessResponses represents a success response with a list of orders.
type SuccessResponses struct {
	Message string  `json:"message"`
	Orders  []Order `json:"orders"`
}

// Customer represents a customer in the system.
type Customer struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Address    string    `json:"address"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	PostalCode string    `json:"postal_code"`
	Country    string    `json:"country"`
	AccountID  uint      `json:"account_id"`
}

// User represents a user in the system.
type User struct {
	ID         uint       `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `gorm:"index" json:"deleted_at" swaggertype:"string" example:"2023-01-01T00:00:00Z"`
	PersonalID string     `json:"personal_id" gorm:"unique;not null"`
	Name       string     `json:"name" gorm:"unique;not null"`
	Email      string     `json:"email" gorm:"unique;not null"`
	Age        int        `json:"age" gorm:"not null"`
	BirthDate  string     `json:"birth_date" gorm:"not null"`
	RoleID     uint       `json:"role_id" gorm:"not null"`
	Role       string     `json:"role" gorm:"foreignKey:RoleID"`
	Permission Permission `json:"permission" gorm:"not null"`
	Phone      string     `json:"phone" gorm:"unique; not null"`
	Street     string     `json:"street"`
	City       string     `json:"city"`
	Password   string     `json:"password" gorm:"not null"`
	IsAdmin    bool       `json:"is_admin" gorm:"default:false"`
	AccountID  uint       `json:"account_id"`
}

// Permission represents user permissions.
type Permission string

const (
	PermissionWorker  Permission = "worker"
	PermissionManager Permission = "manager"
)

// Scan implements the Scanner interface for Permission.
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

// Value implements the Valuer interface for Permission.
func (p Permission) Value() (driver.Value, error) {
	switch p {
	case PermissionWorker, PermissionManager:
		return string(p), nil
	}
	return nil, errors.New("invalid permission value")
}

// String implements the Stringer interface for Permission.
func (p Permission) String() string {
	return string(p)
}

// Role represents a role in the system.
type Role struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at" swaggertype:"string" example:"2023-01-01T00:00:00Z"`
	Role         string     `gorm:"unique;not null"`
	Description  string
	IsActive     bool       `gorm:"default:true"`
	Users        []User     `gorm:"foreignKey:RoleID"`
	DepartmentID uint       `gorm:"not null"`
	Department   Department `gorm:"foreignKey:DepartmentID"`
	AccountID    uint
}

// Department represents a department in the system.
type Department struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at" swaggertype:"string" example:"2023-01-01T00:00:00Z"`
	Name      string     `json:"department"`
	Roles     []Role     `json:"roles"`
	IsActive  bool       `gorm:"default:true"`
	AccountID uint       `json:"account_id"`
}
