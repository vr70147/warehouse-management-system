package handlers

import (
	"encoding/json"
	"net/http"
	"order-processing/internal/cache"
	"order-processing/internal/model"
	"strconv"

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
// @Failure 500 {object} model.ErrorResponse
// @Router /orders [get]
func GetOrders(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		id := c.Query("id")
		if id != "" {
			cachedOrder, err := cache.GetCache(id)
			if err == nil {
				var order model.Order
				json.Unmarshal([]byte(cachedOrder), &order)
				c.JSON(http.StatusOK, order)
				return
			}
		}

		var orders []model.Order
		query := db.Where("account_id = ?", accountID)

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

		if result := query.Find(&orders); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: result.Error.Error()})
			return
		}

		if id != "" {
			orderJSON, _ := json.Marshal(orders)
			cache.SetCache(id, string(orderJSON))
		}

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
func CreateOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		var orderRequest model.Order
		if err := c.ShouldBindJSON(&orderRequest); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid order data"})
			return
		}

		if orderRequest.CustomerID == 0 || orderRequest.Quantity <= 0 || orderRequest.ProductID == 0 {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Missing or invalid fields"})
			return
		}

		tx := db.Begin()
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Database transaction error"})
			return
		}

		order := model.Order{
			AccountID:  accountID.(uint),
			CustomerID: orderRequest.CustomerID,
			Quantity:   orderRequest.Quantity,
			ProductID:  orderRequest.ProductID,
			Status:     "Pending",
		}

		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to create order"})
			return
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to commit transaction"})
			return
		}

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
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		var orderUpdate model.Order
		if err := c.ShouldBindJSON(&orderUpdate); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid order data"})
			return
		}

		var currentOrder model.Order
		if err := db.Where("id = ? AND account_id = ?", orderUpdate.ID, accountID).First(&currentOrder).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Order not found"})
			return
		}

		if orderUpdate.Version != currentOrder.Version {
			c.JSON(http.StatusConflict, model.ErrorResponse{Error: "Order version mismatch"})
			return
		}

		orderUpdate.Version++
		if err := db.Model(&model.Order{}).Where("id = ? AND account_id = ? AND version = ?", orderUpdate.ID, accountID, currentOrder.Version).Updates(map[string]interface{}{
			"status":  orderUpdate.Status,
			"version": orderUpdate.Version,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to update order"})
			return
		}

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
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		id := c.Param("id")
		if result := db.Where("id = ? AND account_id = ?", id, accountID).Delete(&model.Order{}); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to delete order"})
			return
		}
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
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		id := c.Param("id")
		if result := db.Unscoped().Where("id = ? AND account_id = ?", id, accountID).Delete(&model.Order{}); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to delete order"})
			return
		}
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
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		id := c.Param("id")
		var order model.Order

		if result := db.Unscoped().Where("id = ? AND account_id = ?", id, accountID).First(&order); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Order not found"})
			return
		}

		if result := db.Model(&order).Update("DeletedAt", nil); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to recover order"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order recovered successfully"})
	}
}
