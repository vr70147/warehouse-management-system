package handlers

import (
	"net/http"
	"order-processing/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetOrders(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")

		if id != "" {
			// Fetch a single order by ID
			var order model.Order
			if result := db.First(&order, id); result.Error != nil {
				c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Order not found"})
				return
			}
			c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order found", Order: order})
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

		c.JSON(http.StatusOK, model.SuccessResponses{Message: "Orders found", Orders: orders})
	}
}

func CreateOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order model.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		if result := db.Create(&order); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order created successfully", Order: order})
	}
}

func UpdateOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var order model.Order
		if result := db.First(&order, id); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Order not found"})
			return
		}

		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		if result := db.Save(&order); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Order updated successfully", Order: order})
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
