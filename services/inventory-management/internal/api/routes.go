package api

import "github.com/gin-gonic/gin"

func Routers() {
	r := gin.Default()

	r.POST("/inventory/products", CreateProducts)
	r.GET("/inventory/products", GetProducts)
	r.PUT("/inventory/products/:id", UpdateProduct)
	r.DELETE("/inventory/products/:id", DeleteProduct)
	r.DELETE("/inventory/products/hard/:id", DeleteProductPermanently)

	r.Run()
}
