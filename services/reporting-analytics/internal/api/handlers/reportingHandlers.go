package handlers

import (
	"net/http"
	"reporting-analytics/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetSalesReports godoc
// @Summary Get sales reports
// @Description Retrieve sales reports for the account
// @Produce json
// @Success 200 {array} model.SalesReport
// @Failure 401 {object} gin.H{"error": "Account ID not found"}
// @Failure 500 {object} gin.H{"error": "Failed to retrieve sales reports"}
// @Router /reports/sales [get]
func GetSalesReports(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Account ID not found"})
			return
		}

		var reports []model.SalesReport
		if err := db.Where("account_id = ?", accountID).Find(&reports).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sales reports"})
			return
		}

		c.JSON(http.StatusOK, reports)
	}
}

// GetInventoryLevels godoc
// @Summary Get inventory levels
// @Description Retrieve inventory levels for the account
// @Produce json
// @Success 200 {array} model.InventoryLevel
// @Failure 401 {object} gin.H{"error": "Account ID not found"}
// @Failure 500 {object} gin.H{"error": "Failed to retrieve inventory levels"}
// @Router /reports/inventory [get]
func GetInventoryLevels(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Account ID not found"})
			return
		}

		var levels []model.InventoryLevel
		if err := db.Where("account_id = ?", accountID).Find(&levels).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inventory levels"})
			return
		}

		c.JSON(http.StatusOK, levels)
	}
}

// GetShippingStatuses godoc
// @Summary Get shipping statuses
// @Description Retrieve shipping statuses for the account
// @Produce json
// @Success 200 {array} model.ShippingStatus
// @Failure 401 {object} gin.H{"error": "Account ID not found"}
// @Failure 500 {object} gin.H{"error": "Failed to retrieve shipping statuses"}
// @Router /reports/shipping [get]
func GetShippingStatuses(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Account ID not found"})
			return
		}

		var statuses []model.ShippingStatus
		if err := db.Where("account_id = ?", accountID).Find(&statuses).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve shipping statuses"})
			return
		}

		c.JSON(http.StatusOK, statuses)
	}
}

// GetUserActivities godoc
// @Summary Get user activities
// @Description Retrieve user activities for the account
// @Produce json
// @Success 200 {array} model.UserActivity
// @Failure 401 {object} gin.H{"error": "Account ID not found"}
// @Failure 500 {object} gin.H{"error": "Failed to retrieve user activities"}
// @Router /reports/users [get]
func GetUserActivities(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Account ID not found"})
			return
		}

		var activities []model.UserActivity
		if err := db.Where("account_id = ?", accountID).Find(&activities).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user activities"})
			return
		}

		c.JSON(http.StatusOK, activities)
	}
}
