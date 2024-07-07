package routes

import (
	"inventory-management/internal/api/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB) {
	products := r.Group("/products")
	products.POST("/", handlers.CreateProduct(db))
	products.GET("/", handlers.GetProducts(db))
	products.PUT("/:id", handlers.UpdateProduct(db))
	products.DELETE("/:id", handlers.SoftDeleteProduct(db))
	products.DELETE("hard/:id", handlers.HardDeleteProduct(db))
	products.POST("/:id/recover", handlers.RecoverProduct(db))

	categories := r.Group("/categories")
	categories.POST("/", handlers.CreateCategory(db))
	categories.GET("/", handlers.GetCategories(db))
	categories.PUT("/:id", handlers.UpdateCategory(db))
	categories.DELETE("/:id", handlers.SoftDeleteCategory(db))
	categories.POST("/hard/:id", handlers.HardDeleteCategory(db))
	categories.POST("/:id/recover", handlers.RecoverCategory(db))

	stocks := r.Group("/stocks")
	stocks.POST("/", handlers.CreateStock(db))
	stocks.GET("/", handlers.GetStocks(db))
	stocks.PUT("/:id", handlers.UpdateStock(db))
	stocks.DELETE("/:id", handlers.SoftDeleteStock(db))
	stocks.DELETE("/hard/:id", handlers.HardDeleteStock(db))
	stocks.POST("/:id/recover", handlers.RecoverStock(db))

	suppliers := r.Group("/suppliers")
	suppliers.POST("/", handlers.CreateSupplier(db))
	suppliers.GET("/", handlers.GetSuppliers(db))
	suppliers.PUT("/:id", handlers.UpdateSupplier(db))
	suppliers.DELETE("/:id", handlers.SoftDeleteSupplier(db))
	suppliers.DELETE("/hard/:id", handlers.HardDeleteSupplier(db))
	suppliers.POST("/:id/recover", handlers.RecoverSupplier(db))
}
