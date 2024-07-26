package model

import (
	"database/sql/driver"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	AccountID  uint           `gorm:"index" json:"account_id"`
	ProductID  uint           `json:"product_id"`
	Quantity   uint           `json:"quantity"`
	CustomerID uint           `json:"customer_id"`
	Status     string         `json:"status"`
	Version    int            `json:"version"`
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
	Order   Order  `json:"order"`
}

type SuccessResponses struct {
	Message string  `json:"message"`
	Orders  []Order `json:"orders"`
}
