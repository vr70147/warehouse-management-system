package model

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Error string `json:"error"`
}

type UpdateUserResponse struct {
	SuccessResponse
	User User `json:"user"`
}

type UserResponse struct {
	User
	Role       string `json:"role"`
	Permission string `json:"permission"`
	IsActive   bool   `json:"is_active"`
	Department string `json:"department"`
}

type UsersResponse struct {
	Message string         `json:"message"`
	Users   []UserResponse `json:"users"`
}

type DepartmentResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"department"`
	Roles []Role `json:"roles"`
}

type DepartmentsResponse struct {
	Message     string               `json:"message"`
	Departments []DepartmentResponse `json:"department"`
}

type RoleResponse struct {
	Role        Role   `json:"role"`
	Description string `json:"description"`
	Permission  string `json:"permission"`
	IsActive    bool   `json:"is_active"`
	Department  string `json:"department"`
}

// RolesResponse represents the response for retrieving roles
type RolesResponse struct {
	Message string `json:"message"`
	Roles   []Role `json:"roles"`
}
