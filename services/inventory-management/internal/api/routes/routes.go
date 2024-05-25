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
	categories.DELETE("/:id", api.DeleteCategory(db))
	categories.POST("/:id/recover", api.RecoverCategory(db))
	categories.POST("/inventory/categories/hard/:id", api.HardDeleteCategory(db))

	stocks := r.Group("/stocks")
	stocks.POST("/inventory/stocks", api.CreateStock(db))
	stocks.GET("/inventory/stocks", api.GetStocks(db))
	stocks.PUT("/inventory/stocks/:id", api.UpdateStock(db))
	stocks.DELETE("/inventory/stocks/:id", api.DeleteStock(db))
	stocks.DELETE("/inventory/stocks/hard/:id", api.HardDeleteStock(db))
	stocks.POST("/inventory/stocks/:id/recover", api.RecoverStock(db))

	suppliers := r.Group("/suppliers")
	suppliers.POST("/inventory/suppliers", api.CreateSupplier(db))
	suppliers.GET("/inventory/suppliers", api.GetSuppliers(db))
	suppliers.PUT("/inventory/suppliers/:id", api.UpdateSupplier(db))
	suppliers.DELETE("/inventory/suppliers/:id", api.DeleteSupplier(db))
	suppliers.DELETE("/inventory/suppliers/hard/:id", api.HardDeleteSupplier(db))
	suppliers.POST("/inventory/suppliers/:id/recover", api.RecoverSupplier(db))
}
