package model

import (
	"database/sql/driver"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	AccountID    uint           `gorm:"index" json:"account_id"`
	ProductID    uint           `json:"product_id"`
	Quantity     uint           `json:"quantity"`
	CustomerID   uint           `json:"customer_id"`
	Status       string         `json:"status"`
	Version      int            `json:"version"`
	ShippingDate time.Time      `json:"shipping_date"`
}

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

// OrderStatusUpdateRequest godoc
// @Summary Request payload to update the status of an order
// @Description This model is used to capture the new status for updating an order
// @Tags orders
// @Accept json
// @Produce json
type OrderStatusUpdateRequest struct {
	Status string `json:"status" binding:"required"`
}

type OrderEvent struct {
	OrderID   uint   `json:"order_id"`
	ProductID uint   `json:"product_id"`
	Quantity  uint   `json:"quantity"`
	Action    string `json:"action"`
}

type InventoryStatusEvent struct {
	OrderID uint   `json:"order_id"`
	Status  string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Order   Order  `json:"order"`
}
type SuccessCustomerResponse struct {
	Message  string   `json:"message"`
	Order    Order    `json:"order"`
	Customer Customer `json:"customer"`
}

type SuccessResponses struct {
	Message string  `json:"message"`
	Orders  []Order `json:"orders"`
}
