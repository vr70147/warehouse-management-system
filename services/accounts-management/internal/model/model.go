package model

type Account struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Email string `json:"email" gorm:"unique;not null"`
	Name  string `json:"name" gorm:"not null"`
}

type ErrorResponse struct {
	Error string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
