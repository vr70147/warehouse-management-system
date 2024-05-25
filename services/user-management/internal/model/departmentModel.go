package model

import (
	"time"
)

type Department struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"department"`
	Roles     []Role    `json:"roles"`
}
