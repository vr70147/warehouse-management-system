package routes

import (
	"order-processing/internal/api/handlers"
	"order-processing/internal/middleware"
	"order-processing/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB, ns *utils.NotificationService) {
	r.Use(middleware.CORSMiddleware())

	orders := r.Group("/orders")
	orders.Use(middleware.AuthMiddleware(db))

	orders.GET("", handlers.GetOrders(db))
	orders.POST("", handlers.CreateOrder(db, ns))
	orders.PUT("/:id", handlers.UpdateOrder(db))
	orders.DELETE("/:id", handlers.SoftDeleteOrder(db))
	orders.DELETE("/hard/:id", handlers.HardDeleteOrder(db))
	orders.POST("/recover/:id", handlers.RecoverOrder(db))
	orders.POST("/cancel/:id", handlers.CancelOrder(db, ns))
	orders.PUT("/:id/status", handlers.UpdateOrderStatus(db))
}
