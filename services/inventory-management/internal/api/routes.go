package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// router function
func Routers(r *gin.Engine, db *gorm.DB) {
	products := r.Group("/products")

	products.POST("/inventory/products", createProduct(db))
	products.GET("/inventory/products", getProducts(db))
	products.PUT("/inventory/products/:id", updateProduct(db))
	products.DELETE("/inventory/products/:id", softDeleteProduct(db))
	products.DELETE("/inventory/products/hard/:id", hardDeleteProduct(db))

}
