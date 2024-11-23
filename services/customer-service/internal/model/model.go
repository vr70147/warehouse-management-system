package model

import "time"

type Customer struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Name       string    `json:"name" gorm:"not null"`
	Email      string    `json:"email" gorm:"unique;not null"`
	Phone      string    `json:"phone" gorm:"unique;not null"`
	Address    string    `json:"address"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	PostalCode string    `json:"postal_code"`
	Country    string    `json:"country"`
	AccountID  uint      `json:"account_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message  string   `json:"message"`
	Customer Customer `json:"customer"`
}

type SuccessResponses struct {
	Message   string     `json:"message"`
	Customers []Customer `json:"customer"`
}
