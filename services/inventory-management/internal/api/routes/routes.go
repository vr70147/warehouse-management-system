package routes

import (
	"inventory-management/internal/api"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB) {
	products := r.Group("/products")
	products.POST("/", api.CreateProduct(db))
	products.GET("/", api.GetProducts(db))
	products.PUT("/:id", api.UpdateProduct(db))
	products.DELETE("/:id", api.SoftDeleteProduct(db))
	products.DELETE("hard/:id", api.HardDeleteProduct(db))
	products.POST("/:id/recover", api.RecoverProduct(db))

	categories := r.Group("/categories")
	categories.POST("/", api.CreateCategory(db))
	categories.GET("/", api.GetCategories(db))
	categories.PUT("/:id", api.UpdateCategory(db))
	categories.DELETE("/:id", api.SoftDeleteCategory(db))
	categories.POST("/hard/:id", api.HardDeleteCategory(db))
	categories.POST("/:id/recover", api.RecoverCategory(db))

	stocks := r.Group("/stocks")
	stocks.POST("/", api.CreateStock(db))
	stocks.GET("/", api.GetStocks(db))
	stocks.PUT("/:id", api.UpdateStock(db))
	stocks.DELETE("/:id", api.SoftDeleteStock(db))
	stocks.DELETE("/hard/:id", api.HardDeleteStock(db))
	stocks.POST("/:id/recover", api.RecoverStock(db))

	suppliers := r.Group("/suppliers")
	suppliers.POST("/", api.CreateSupplier(db))
	suppliers.GET("/", api.GetSuppliers(db))
	suppliers.PUT("/:id", api.UpdateSupplier(db))
	suppliers.DELETE("/:id", api.SoftDeleteSupplier(db))
	suppliers.DELETE("/hard/:id", api.HardDeleteSupplier(db))
	suppliers.POST("/:id/recover", api.RecoverSupplier(db))
}
