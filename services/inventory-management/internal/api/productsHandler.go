package api

import (
	"inventory-management/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Create a new product
// @Description Add a new product to the inventory
// @Tags products
// @Accept json
// @Produce json
// @Param product body model.Product true "Product to create"
// @Success 200 {object} model.Product
// @Failure 400 {object} model.ErrorResponse
// @Router /products [post]
func CreateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product model.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}
		db.Create(&product)
		c.JSON(http.StatusOK, product)
	}
}

// @Summary Get all products or filter by various criteria
// @Description Retrieve all products or filter by ID, name, category ID, or supplier ID
// @Tags products
// @Produce json
// @Param id query int false "Product ID"
// @Param name query string false "Product name"
// @Param category_id query int false "Category ID"
// @Param supplier_id query int false "Supplier ID"
// @Success 200 {array} model.Product
// @Router /products [get]
func GetProducts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []model.Product
		id := c.Query("id")
		name := c.Query("name")
		categoryID := c.Query("category_id")
		supplierID := c.Query("supplier_id")

		query := db.Preload("Category").Preload("Supplier").Preload("Stocks")

		if id != "" {
			if err := query.Where("id = ?", id).First(&products).Error; err != nil {
				c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Product not found"})
				return
			}
			c.JSON(http.StatusOK, products[0])
			return
		}

		if name != "" {
			query = query.Where("name LIKE ?", "%"+name+"%")
		}

		if categoryID != "" {
			query = query.Where("category_id = ?", categoryID)
		}

		if supplierID != "" {
			query = query.Where("supplier_id = ?", supplierID)
		}

		query.Find(&products)
		c.JSON(http.StatusOK, products)
	}
}

// @Summary Update an existing product
// @Description Update details of an existing product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body model.Product true "Product to update"
// @Success 200 {object} model.Product
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /products/{id} [put]
func UpdateProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product model.Product
		if err := db.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Product not found"})
			return
		}

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		db.Save(&product)
		c.JSON(http.StatusOK, product)
	}
}

// @Summary Soft delete an existing product
// @Description Soft delete a product by ID
// @Tags products
// @Param id path int true "Product ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /products/{id} [delete]
func SoftDeleteProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Where("id = ?", c.Param("id")).Delete(&model.Product{}).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Product not found"})
			return
		}
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Product soft deleted"})
	}
}

// @Summary Hard delete an existing product and its stocks
// @Description Hard delete a product by ID along with its associated stocks
// @Tags products
// @Param id path int true "Product ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /products/hard/{id} [delete]
func HardDeleteProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.Begin()

		// Delete associated stocks
		if err := tx.Where("product_id = ?", c.Param("id")).Unscoped().Delete(&model.Stock{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to delete associated stocks"})
			return
		}

		// Delete the product
		if err := tx.Unscoped().Where("id = ?", c.Param("id")).Delete(&model.Product{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Product not found"})
			return
		}

		tx.Commit()
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Product and associated stocks hard deleted"})
	}
}

// @Summary Recover a soft-deleted product
// @Description Recover a soft-deleted product by ID
// @Tags products
// @Param id path int true "Product ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /products/{id}/recover [post]
func RecoverProduct(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productID := c.Param("id")

		if err := db.Unscoped().Model(&model.Product{}).Where("id = ?", productID).Update("deleted_at", nil).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to recover product",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Product recovered",
		})
	}
}
