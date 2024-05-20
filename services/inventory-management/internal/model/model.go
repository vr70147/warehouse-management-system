package model

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `json:"category"`
	SupplierID  uint     `json:"supplier_id"`
	Supplier    Supplier `json:"supplier"`
	Stocks      []Stock  `json:"stocks" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
}

type Stock struct {
	gorm.Model
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  uint    `json:"quantity"`
	Location  string  `json:"location"`
}

type Category struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Products    []Product `json:"products"`
}

type Supplier struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Products    []Product `json:"products"`
}
