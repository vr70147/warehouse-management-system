package handlers

import (
	"net/http"
	"reporting-analytics/internal/initializers"
	"reporting-analytics/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetSalesReports(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reports []model.SalesReport
		initializers.DB.Find(&reports)
		c.JSON(http.StatusOK, reports)
	}
}

func GetInventoryLevels(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var levels []model.InventoryLevel
		initializers.DB.Find(&levels)
		c.JSON(http.StatusOK, levels)
	}
}

func GetShippingStatuses(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var statuses []model.ShippingStatus
		initializers.DB.Find(&statuses)
		c.JSON(http.StatusOK, statuses)
	}
}

func GetUserActivities(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var activities []model.UserActivity
		initializers.DB.Find(&activities)
		c.JSON(http.StatusOK, activities)
	}
}
