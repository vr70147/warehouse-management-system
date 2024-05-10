package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email" gorm:"unique"`
	Age       int    `json:"age"`
	BirthDate string `json:"birthDate"`
	Role      string `json:"role"`
	Phone     string `json:"phone" gorm:"unique"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Password  string `json:"password"`
}
