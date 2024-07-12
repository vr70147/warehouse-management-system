package routes

import (
	"reporting-analytics/internal/api/handlers"
	"reporting-analytics/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB) {

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware(db))

	report := r.Group("/reports")

	report.GET("/sales", handlers.GetSalesReports(db))
	report.GET("/inventory", handlers.GetInventoryLevels(db))
	report.GET("/shipping", handlers.GetShippingStatuses(db))
	report.GET("/user-activity", handlers.GetUserActivities(db))
}
