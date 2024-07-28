package routes

import (
	"customer-service/internal/api/handlers"
	"customer-service/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB) {
	r.Use(middleware.CORSMiddleware())

	customers := r.Group("/customers")
	customers.Use(middleware.AuthMiddleware(db))

	customers.GET("", handlers.GetCustomers(db))
	customers.GET("/:id", handlers.GetCustomer(db))
	customers.POST("", handlers.CreateCustomer(db))
	customers.PUT("/:id", handlers.UpdateCustomer(db))
	customers.DELETE("/:id", handlers.SoftDeleteCustomer(db))
	customers.DELETE("/hard/:id", handlers.HardDeleteCustomer(db))
	customers.POST("/recover/:id", handlers.RecoverCustomer(db))
}
