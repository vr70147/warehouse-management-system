package model

type OrderRequest struct {
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
	Quantity  uint   `json:"quantity"`
	Address   string `json:"address"`
}

type OrderResponse struct {
	OrderID uint   `json:"order_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type InventoryRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}

type InventoryResponse struct {
	ProductID uint   `json:"product_id"`
	Quantity  uint   `json:"quantity"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type ShippingRequest struct {
	OrderID uint   `json:"order_id"`
	Address string `json:"address"`
}

type ShippingResponse struct {
	OrderID uint   `json:"order_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type SalesReportRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type SalesReportResponse struct {
	OrderID    uint    `json:"order_id"`
	ProductID  uint    `json:"product_id"`
	Quantity   uint    `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	Timestamp  string  `json:"timestamp"`
}
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
