package model

import (
	"time"
)

type Department struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Name      string     `json:"department"`
	Roles     []Role     `json:"roles"`
	IsActive  bool       `gorm:"default:true"`
	AccountID uint       `json:"account_id"`
}
