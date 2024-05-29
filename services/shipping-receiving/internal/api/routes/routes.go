package routes

import (
	"shipping-receiving/internal/api"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB) {
	shippings := r.Group("/shipping-receiving")
	shippings.POST("/", api.CreateShipping(db))
	shippings.GET("/", api.GetShippings(db))
	shippings.PUT("/:id", api.UpdateShipping(db))
	shippings.DELETE("/:id", api.SoftDeleteShipping(db))
	shippings.DELETE("/hard/:id", api.HardDeleteShipping(db))
	shippings.POST("/:id/recover", api.RecoverShipping(db))
}
