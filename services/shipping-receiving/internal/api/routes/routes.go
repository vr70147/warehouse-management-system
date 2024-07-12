package routes

import (
	"shipping-receiving/internal/api/handlers"
	"shipping-receiving/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB) {
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware(db))

	shippings := r.Group("/shipping-receiving")
	shippings.POST("/", handlers.CreateShipping(db))
	shippings.GET("/", handlers.GetShippings(db))
	shippings.PUT("/:id", handlers.UpdateShipping(db))
	shippings.DELETE("/:id", handlers.SoftDeleteShipping(db))
	shippings.DELETE("/hard/:id", handlers.HardDeleteShipping(db))
	shippings.POST("/:id/recover", handlers.RecoverShipping(db))
}
