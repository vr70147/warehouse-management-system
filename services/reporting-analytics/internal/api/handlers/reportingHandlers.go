package handlers

import (
	"net/http"
	"reporting-analytics/internal/initializers"
	"reporting-analytics/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetSalesReports godoc
// @Summary Get sales reports
// @Description Get sales reports
// @Produce json
// @Success 200 {array} model.SalesReport
// @Router /reports/sales [get]
func GetSalesReports(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reports []model.SalesReport
		initializers.DB.Find(&reports)
		c.JSON(http.StatusOK, reports)
	}
}

// GetInventoryLevels godoc
// @Summary Get inventory levels
// @Description Get inventory levels
// @Produce json
// @Success 200 {array} model.InventoryLevel
// @Router /reports/inventory [get]
func GetInventoryLevels(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var levels []model.InventoryLevel
		initializers.DB.Find(&levels)
		c.JSON(http.StatusOK, levels)
	}
}

// GetShippingStatuses godoc
// @Summary Get shipping statuses
// @Description Get shipping statuses
// @Produce json
// @Success 200 {array} model.ShippingStatus
// @Router /reports/shipping [get]
func GetShippingStatuses(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statuses []model.ShippingStatus
		initializers.DB.Find(&statuses)
		c.JSON(http.StatusOK, statuses)
	}
}

// GetUserActivities godoc
// @Summary Get user activities
// @Description Get user activities
// @Produce json
// @Success 200 {array} model.UserActivity
// @Router /reports/users [get]
func GetUserActivities(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var activities []model.UserActivity
		initializers.DB.Find(&activities)
		c.JSON(http.StatusOK, activities)
	}
}
