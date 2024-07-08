package handlers

import (
	"inventory-management/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateStock godoc
// @Summary Create a new stock item
// @Description Create a new stock item in the inventory
// @Tags stocks
// @Accept json
// @Produce json
// @Param body body model.Stock true "Stock data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /stocks [post]
func CreateStock(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var stock model.Stock

		if err := c.ShouldBindJSON(&stock); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Failed to read body",
			})
			return
		}

		if result := db.Create(&stock); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to create stock",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Stock created successfully",
		})
	}
}

// UpdateStock godoc
// @Summary Update a stock item
// @Description Update a stock item by ID
// @Tags stocks
// @Accept json
// @Produce json
// @Param id path int true "Stock ID"
// @Param body body model.Stock true "Stock data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /stocks/{id} [put]
func UpdateStock(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		stockID := c.Param("id")
		var stock model.Stock

		if result := db.First(&stock, stockID); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "Stock not found",
			})
			return
		}

		if err := c.ShouldBindJSON(&stock); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Invalid request data",
			})
			return
		}

		if result := db.Save(&stock); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to update stock",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Stock updated successfully",
		})
	}
}

// GetStocks godoc
// @Summary Get all stock items
// @Description Retrieve all stock items or filter by query parameters
// @Tags stocks
// @Produce json
// @Success 200 {object} model.StocksResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /stocks [get]
func GetStocks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var stocks []model.Stock
		if result := db.Find(&stocks); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to retrieve stocks",
			})
			return
		}

		var stockResponses []model.StockResponse
		for _, stock := range stocks {
			stockResponses = append(stockResponses, model.StockResponse{
				ID:          stock.ID,
				ProductName: stock.Product.Name,
				Quantity:    int(stock.Quantity),
				Location:    stock.Location,
			})
		}

		c.JSON(http.StatusOK, model.StocksResponse{
			Message: "Stocks retrieved successfully",
			Stocks:  stockResponses,
		})
	}
}

// DeleteStock godoc
// @Summary Delete a stock item
// @Description Delete a stock item by ID
// @Tags stocks
// @Produce json
// @Param id path int true "Stock ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /stocks/{id} [delete]
func SoftDeleteStock(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		stockID := c.Param("id")

		if result := db.Delete(&model.Stock{}, stockID); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to delete stock",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Stock deleted successfully",
		})
	}
}

// HardDeleteStock godoc
// @Summary Hard delete a stock item
// @Description Permanently delete a stock item by ID
// @Tags stocks
// @Produce json
// @Param id path int true "Stock ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /stocks/{id}/hard [delete]
func HardDeleteStock(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		stockID := c.Param("id")

		if result := db.Unscoped().Delete(&model.Stock{}, stockID); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to delete stock",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Stock deleted permanently",
		})
	}
}

// RecoverStock godoc
// @Summary Recover a deleted stock item
// @Description Recover a soft-deleted stock item by ID
// @Tags stocks
// @Produce json
// @Param id path int true "Stock ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /stocks/{id}/recover [post]
func RecoverStock(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		stockID := c.Param("id")

		if result := db.Model(&model.Stock{}).Unscoped().Where("id = ?", stockID).Update("deleted_at", nil); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to recover stock",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Stock recovered successfully",
		})
	}
}