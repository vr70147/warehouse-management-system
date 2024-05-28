package routes

import (
	"order-processing/internal/api"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB) {
	orders := r.Group("/orders")
	orders.POST("/", api.CreateOrder(db))
	orders.GET("/", api.GetOrders(db))
	orders.PUT("/:id", api.UpdateOrder(db))
	orders.DELETE("/:id", api.SoftDeleteOrder(db))
	orders.DELETE("hard/:id", api.HardDeleteOrder(db))
	orders.POST("/:id/recover", api.RecoverOrder(db))
}
