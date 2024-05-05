package model

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Role     string `json:"role"`
	Phone    string `json:"phone"`
	Street   string `json:"street"`
	City     string `json:"city"`
	Password string `json:"password"`
}
