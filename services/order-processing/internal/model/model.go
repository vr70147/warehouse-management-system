package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	ProductID  uint
	Quantity   uint
	CustomerID uint
	Status     string
}
