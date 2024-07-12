package model

type OrderRequest struct {
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
	Quantity  uint   `json:"quantity"`
	Address   string `json:"address"`
	AccountID uint   `json:"account_id"`
}

type OrderResponse struct {
	OrderID   uint   `json:"order_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	AccountID uint   `json:"account_id"`
}

type InventoryRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
	AccountID uint `json:"account_id"`
}

type InventoryResponse struct {
	ProductID uint   `json:"product_id"`
	Quantity  uint   `json:"quantity"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	AccountID uint   `json:"account_id"`
}

type ShippingRequest struct {
	OrderID   uint   `json:"order_id"`
	Address   string `json:"address"`
	AccountID uint   `json:"account_id"`
}

type ShippingResponse struct {
	OrderID   uint   `json:"order_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	AccountID uint   `json:"account_id"`
}

type UserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AccountID uint   `json:"account_id"`
}

type UserResponse struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	AccountID uint   `json:"account_id"`
}

type SalesReportRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	AccountID uint   `json:"account_id"`
}

type SalesReportResponse struct {
	OrderID    uint    `json:"order_id"`
	ProductID  uint    `json:"product_id"`
	Quantity   uint    `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	Timestamp  string  `json:"timestamp"`
	AccountID  uint    `json:"account_id"`
}
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
