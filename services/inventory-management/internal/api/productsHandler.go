package api

import (
	"inventory-management/internal/initializers"
	"inventory-management/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateProducts(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindBodyWithJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	initializers.DB.Create(&product)
	c.JSON(http.StatusOK, gin.H{
		"message": "Product created successfuly",
	})
}

func GetProducts(c *gin.Context) {
	var products []model.Product
	id := c.Query("id")
	name := c.Query("name")
	categoryID := c.Query("category")
	supplierID := c.Query("supplier")

	query := initializers.DB.Preload("Category").Preload("Supplier").Preload("Stocks")

	if id != "" {
		if err := query.Where("id = ?", id).First(&products).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}
		c.JSON(http.StatusOK, products[0])
		return
	}

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if categoryID != "" {
		query = query.Where("category = ?", "%"+categoryID+"%")
	}

	if supplierID != "" {
		query = query.Where("supplier = ?", "%"+supplierID+"%")
	}

	query.Find(&products)
	c.JSON(http.StatusOK, products)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product model.Product
	if err := initializers.DB.Where("id = ?", id).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	initializers.DB.Save(&product)
	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := initializers.DB.Where("id = ?", id).Delete(&model.Product{}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfuly"})
}

func DeleteProductPermanently(c *gin.Context) {
	id := c.Param("id")
	trns := initializers.DB.Begin()

	if err := trns.Where("product_id = ?", id).Unscoped().Delete(&model.Stock{}).Error; err != nil {
		trns.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete product",
		})
		return
	}

	if err := trns.Unscoped().Where("id = ?", c.Param("id")).Delete(&model.Product{}).Error; err != nil {
		trns.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	trns.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Product and associated stocks hard deleted"})
}
