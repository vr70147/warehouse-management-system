package handlers

import (
	"encoding/json"
	"net/http"
	"order-processing/internal/cache"
	"order-processing/internal/initializers"
	"order-processing/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetOrders(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")

		cachedOrder, err := cache.GetCache(id)
		if err == nil {
			var order model.Order
			json.Unmarshal([]byte(cachedOrder), &order)
			c.JSON(http.StatusOK, order)
			return
		}

		// Fetch list of orders with optional query parameters
		var orders []model.Order
		query := db

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

		// Cache the list of orders
		orderJSON, _ := json.Marshal(orders)
		cache.SetCache(id, string(orderJSON))

		c.JSON(http.StatusOK, model.SuccessResponses{Message: "Orders found", Orders: orders})
	}
}

func CreateOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var orderRequest model.Order
		if err := c.ShouldBindJSON(&orderRequest); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		tx := initializers.DB.Begin()
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: tx.Error.Error()})
			return
		}

		order := model.Order{
			CustomerID: orderRequest.CustomerID,
			Quantity:   orderRequest.Quantity,
			ProductID:  orderRequest.ProductID,
			Status:     "Pending",
		}

		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order created successfully", Order: order})
	}
}

func UpdateOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var orderUpdate model.Order

		if err := c.ShouldBindJSON(&orderUpdate); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		var currentOrder model.Order
		if err := initializers.DB.Where("id = ?", orderUpdate.ID).First(&currentOrder).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Order not found"})
			return
		}

		if orderUpdate.Version != currentOrder.Version {
			c.JSON(http.StatusConflict, model.ErrorResponse{Error: "Order version mismatch"})
			return
		}

		orderUpdate.Version++
		if err := initializers.DB.Model(&model.Order{}).Where("id = ? AND version = ?", orderUpdate.ID, currentOrder.Version).Updates(map[string]interface{}{
			"status":  orderUpdate.Status,
			"version": orderUpdate.Version,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order updated successfully", Order: orderUpdate})
	}
}

func SoftDeleteOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if result := db.Delete(&model.Order{}, id); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order deleted successfully"})
	}
}

func HardDeleteOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if result := db.Unscoped().Delete(&model.Order{}, id); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order deleted permanently"})
	}
}

func RecoverOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var order model.Order

		// Find the soft-deleted order
		if result := db.Unscoped().First(&order, id); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Order not found"})
			return
		}

		// Recover the order by setting DeletedAt to NULL
		if result := db.Model(&order).Update("DeletedAt", nil); result.Error != nil {

			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order recovered successfully"})
	}
}
