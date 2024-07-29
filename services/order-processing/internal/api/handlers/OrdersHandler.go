package handlers

import (
	"encoding/json"
	"net/http"
	"order-processing/internal/cache"
	"order-processing/internal/kafka"
	"order-processing/internal/model"
	"order-processing/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetOrders godoc
// @Summary Get orders
// @Description Retrieve a list of orders with optional query parameters
// @Tags orders
// @Produce json
// @Param id query string false "Order ID"
// @Param status query string false "Order Status"
// @Param customer_id query string false "Customer ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} model.SuccessResponses
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /orders [get]
func GetOrders(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve account ID from context
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		// Check for cached order by ID
		id := c.Query("id")
		if id != "" {
			cachedOrder, err := cache.GetCache(id)
			if err == nil {
				var order model.Order
				json.Unmarshal([]byte(cachedOrder), &order)
				c.JSON(http.StatusOK, model.SuccessResponses{Message: "Orders found", Orders: []model.Order{order}})
				return
			}
		}

		var orders []model.Order
		query := db.Where("account_id = ?", accountID)

		// Apply filters based on query parameters
		if status := c.Query("status"); status != "" {
			query = query.Where("status = ?", status)
		}

		if customerID := c.Query("customer_id"); customerID != "" {
			query = query.Where("customer_id = ?", customerID)
		}

		if limit := c.Query("limit"); limit != "" {
			if limitInt, err := strconv.Atoi(limit); err == nil {
				query = query.Limit(limitInt)
			}
		}

		if offset := c.Query("offset"); offset != "" {
			if offsetInt, err := strconv.Atoi(offset); err == nil {
				query = query.Offset(offsetInt)
			}
		}

		// Retrieve orders from the database
		if result := query.Find(&orders); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: result.Error.Error()})
			return
		}

		// Cache the retrieved order by ID
		if id != "" {
			orderJSON, _ := json.Marshal(orders)
			cache.SetCache(id, string(orderJSON))
		}

		// Respond with the retrieved orders
		c.JSON(http.StatusOK, model.SuccessResponses{Message: "Orders found", Orders: orders})
	}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order in the system
// @Tags orders
// @Accept json
// @Produce json
// @Param body body model.Order true "Order data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /orders [post]
func CreateOrder(db *gorm.DB, ns *utils.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve account ID from context
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		// Bind JSON to Order Request
		var orderRequest model.Order
		if err := c.ShouldBindJSON(&orderRequest); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid order data"})
			return
		}

		// Validate Order Fields
		if orderRequest.CustomerID == 0 || orderRequest.Quantity <= 0 || orderRequest.ProductID == 0 {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Missing or invalid fields"})
			return
		}

		// Begin Database Transaction
		tx := db.Begin()
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Database transaction error"})
			return
		}

		// Create the order with status "Pending"
		order := model.Order{
			AccountID:  accountID.(uint),
			CustomerID: orderRequest.CustomerID,
			Quantity:   orderRequest.Quantity,
			ProductID:  orderRequest.ProductID,
			Status:     "Pending",
		}

		// Save the new order to the database
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to create order"})
			return
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to commit transaction"})
			return
		}

		// Publish Kafka Event
		kafka.PublishOrderEvent(order.ID, order.ProductID, order.Quantity, "create")

		// Respond with success message
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order created successfully", Order: order})
	}
}

// UpdateOrder godoc
// @Summary Update an order
// @Description Update an order by ID
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param body body model.Order true "Order data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /orders/{id} [put]
func UpdateOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve account ID from context
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		// Bind JSON to Order Update Request
		var orderUpdate model.Order
		if err := c.ShouldBindJSON(&orderUpdate); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid order data"})
			return
		}

		// Retrieve the current order from the database
		var currentOrder model.Order
		if err := db.Where("id = ? AND account_id = ?", orderUpdate.ID, accountID).First(&currentOrder).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Order not found"})
			return
		}

		// Check for version mismatch
		if orderUpdate.Version != currentOrder.Version {
			c.JSON(http.StatusConflict, model.ErrorResponse{Error: "Order version mismatch"})
			return
		}

		// Update the order status and version
		orderUpdate.Version++
		if err := db.Model(&model.Order{}).Where("id = ? AND account_id = ? AND version = ?", orderUpdate.ID, accountID, currentOrder.Version).Updates(map[string]interface{}{
			"status":  orderUpdate.Status,
			"version": orderUpdate.Version,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to update order"})
			return
		}

		// Respond with success message
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order updated successfully", Order: orderUpdate})
	}
}

// SoftDeleteOrder godoc
// @Summary Soft delete an order
// @Description Soft delete an order by ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /orders/{id} [delete]
func SoftDeleteOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve account ID from context
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		// Retrieve the order ID from the path
		id := c.Param("id")
		if result := db.Where("id = ? AND account_id = ?", id, accountID).Delete(&model.Order{}); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to delete order"})
			return
		}

		// Respond with success message
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order deleted successfully"})
	}
}

// HardDeleteOrder godoc
// @Summary Hard delete an order
// @Description Permanently delete an order by ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /orders/{id}/hard [delete]
func HardDeleteOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve account ID from context
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		// Retrieve the order ID from the path
		id := c.Param("id")
		if result := db.Unscoped().Where("id = ? AND account_id = ?", id, accountID).Delete(&model.Order{}); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to delete order"})
			return
		}

		// Respond with success message
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order deleted permanently"})
	}
}

// RecoverOrder godoc
// @Summary Recover a deleted order
// @Description Recover a soft-deleted order by ID
// @Tags orders
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /orders/{id}/recover [post]
func RecoverOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve account ID from context
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		// Retrieve the order ID from the path
		id := c.Param("id")
		var order model.Order

		// Retrieve the soft-deleted order
		if result := db.Unscoped().Where("id = ? AND account_id = ?", id, accountID).First(&order); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Order not found"})
			return
		}

		// Recover the order by updating the DeletedAt field to nil
		if result := db.Model(&order).Update("DeletedAt", nil); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to recover order"})
			return
		}

		// Respond with success message
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order recovered successfully"})
	}
}

// CancelOrder godoc
// @Summary Cancel an order
// @Description Mark an order as cancelled and send a notification email
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /orders/cancel/{id} [post]
func CancelOrder(db *gorm.DB, ns *utils.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve account ID from context
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		// Retrieve the order ID from the path
		orderID := c.Param("id")

		// Retrieve the order from the database
		var order model.Order
		if err := db.Where("id = ? AND account_id = ?", orderID, accountID).First(&order).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Order not found"})
			return
		}

		// Update the order status to "cancelled"
		order.Status = "cancelled"
		if result := db.Save(&order); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to update order"})
			return
		}

		// Send email notification (commented out for now)
		// go func(email string, orderID uint) {
		// 	err := ns.SendOrderCancellationNotification(email, orderID)
		// 	if err != nil {
		// 		// Log the error if the email fails to send
		// 		log.Printf("Failed to send order cancellation notification: %v", err)
		// 	}
		// }(order.CustomerEmail, order.ID)

		// Publish order cancellation event to Kafka
		kafka.PublishOrderEvent(order.ID, order.ProductID, order.Quantity, "cancelled")

		// Respond with success message
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order cancelled successfully"})
	}
}

// UpdateOrderStatus godoc
// @Summary Update the status and shipping date of an order
// @Description Update the status and shipping date of an order
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param status body model.OrderStatusUpdate true "Order Status Update"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /orders/{id}/status [put]
func UpdateOrderStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind JSON to status update request
		var input struct {
			Status       string    `json:"status"`
			ShippingDate time.Time `json:"shipping_date"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request data"})
			return
		}

		// Retrieve the order ID from the path
		orderID := c.Param("id")
		var order model.Order

		// Retrieve the order from the database
		if err := db.First(&order, orderID).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Order not found"})
			return
		}

		// Update the order status and shipping date
		order.Status = input.Status
		order.ShippingDate = input.ShippingDate

		// Save the updated order to the database
		if err := db.Save(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to update order"})
			return
		}

		// Respond with success message
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order updated successfully", Order: order})
	}
}
