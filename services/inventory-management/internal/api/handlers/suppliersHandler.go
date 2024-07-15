package handlers

import (
	"inventory-management/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateSupplier godoc
// @Summary Create a new supplier
// @Description Create a new supplier in the inventory
// @Tags suppliers
// @Accept json
// @Produce json
// @Param body body model.Supplier true "Supplier data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /suppliers [post]
func CreateSupplier(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		var supplier model.Supplier
		if err := c.ShouldBindJSON(&supplier); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request data"})
			return
		}

		supplier.AccountID = accountID.(uint)
		if result := db.Create(&supplier); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to create supplier"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Supplier created successfully"})
	}
}

// UpdateSupplier godoc
// @Summary Update a supplier
// @Description Update a supplier by ID
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path int true "Supplier ID"
// @Param body body model.Supplier true "Supplier data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /suppliers/{id} [put]
func UpdateSupplier(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		supplierID := c.Param("id")
		var supplier model.Supplier
		if result := db.Where("id = ? AND account_id = ?", supplierID, accountID).First(&supplier); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Supplier not found"})
			return
		}

		if err := c.ShouldBindJSON(&supplier); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request data"})
			return
		}

		supplier.AccountID = accountID.(uint)
		if result := db.Save(&supplier); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to update supplier"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Supplier updated successfully"})
	}
}

// GetSuppliers godoc
// @Summary Get all suppliers
// @Description Retrieve all suppliers
// @Tags suppliers
// @Produce json
// @Success 200 {object} model.SuppliersResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /suppliers [get]
func GetSuppliers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		var suppliers []model.Supplier
		if result := db.Where("account_id = ?", accountID).Find(&suppliers); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to retrieve suppliers"})
			return
		}

		var supplierResponses []model.SupplierResponse
		for _, supplier := range suppliers {
			supplierResponses = append(supplierResponses, model.SupplierResponse{
				ID:          supplier.ID,
				Name:        supplier.Name,
				Description: supplier.Description,
				Email:       supplier.Email,
				Contact:     supplier.Contact,
			})
		}

		c.JSON(http.StatusOK, model.SuppliersResponse{
			Message:   "Suppliers retrieved successfully",
			Suppliers: supplierResponses,
		})
	}
}

// SoftDeleteSupplier godoc
// @Summary Soft delete a supplier
// @Description Soft deletes a supplier by ID and sets the supplier field in related products to null
// @Tags suppliers
// @Produce json
// @Param id path int true "Supplier ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /suppliers/{id} [delete]
func SoftDeleteSupplier(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		supplierID := c.Param("id")

		// Find the supplier by ID
		var supplier model.Supplier
		if err := db.Where("id = ? AND account_id = ?", supplierID, accountID).First(&supplier).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Supplier not found"})
			return
		}

		// Set supplier field in related products to null
		if err := db.Model(&model.Product{}).Where("supplier_id = ? AND account_id = ?", supplierID, accountID).Update("supplier_id", nil).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to reassign products"})
			return
		}

		// Soft delete the supplier
		if err := db.Delete(&supplier).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to delete supplier"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Supplier deleted successfully"})
	}
}

// HardDeleteSupplier godoc
// @Summary Hard delete a supplier
// @Description Hard deletes a supplier by ID and sets the supplier field in related products to null
// @Tags suppliers
// @Produce json
// @Param id path int true "Supplier ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /suppliers/hard/{id} [delete]
func HardDeleteSupplier(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		supplierID := c.Param("id")

		// Find the supplier by ID
		var supplier model.Supplier
		if err := db.Unscoped().Where("id = ? AND account_id = ?", supplierID, accountID).First(&supplier).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Supplier not found"})
			return
		}

		// Set supplier field in related products to null
		if err := db.Model(&model.Product{}).Where("supplier_id = ? AND account_id = ?", supplierID, accountID).Update("supplier_id", nil).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to reassign products"})
			return
		}

		// Hard delete the supplier
		if err := db.Unscoped().Delete(&supplier).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to delete supplier"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Supplier deleted permanently"})
	}
}

// RecoverSupplier godoc
// @Summary Recover a deleted supplier
// @Description Recover a soft-deleted supplier by ID and reassign its products back to the supplier
// @Tags suppliers
// @Produce json
// @Param id path int true "Supplier ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /suppliers/{id}/recover [post]
func RecoverSupplier(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		supplierID := c.Param("id")

		// Recover the soft-deleted supplier by setting deleted_at to NULL
		if result := db.Model(&model.Supplier{}).Unscoped().Where("id = ? AND account_id = ?", supplierID, accountID).Update("deleted_at", nil); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to recover supplier"})
			return
		}

		// Reassign products back to the recovered supplier
		if err := db.Model(&model.Product{}).Where("supplier_id IS NULL AND account_id = ? AND previous_supplier_id = ?", accountID, supplierID).Update("supplier_id", supplierID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to reassign products"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Supplier recovered successfully"})
	}
}
