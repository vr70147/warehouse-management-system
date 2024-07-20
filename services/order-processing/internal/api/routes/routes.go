package routes

import (
	"order-processing/internal/api/handlers"
	"order-processing/internal/api/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB) {

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware(db))

	orders := r.Group("/orders")
	orders.POST("/", handlers.CreateOrder(db))
	orders.GET("/", handlers.GetOrders(db))
	orders.PUT("/:id", handlers.UpdateOrder(db))
	orders.DELETE("/:id", handlers.SoftDeleteOrder(db))
	orders.DELETE("hard/:id", handlers.HardDeleteOrder(db))
	orders.POST("/:id/recover", handlers.RecoverOrder(db))
	orders.POST("/orders/cancel/:id", handlers.CancelOrder(db))
}
