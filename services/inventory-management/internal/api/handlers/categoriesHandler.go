package handlers

import (
	"inventory-management/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category in the inventory
// @Tags categories
// @Accept json
// @Produce json
// @Param body body model.Category true "Category data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /categories [post]
func CreateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		var category model.Category
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request data"})
			return
		}

		category.AccountID = accountID.(uint)
		if result := db.Create(&category); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to create category"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Category created successfully"})
	}
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param body body model.Category true "Category data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /categories/{id} [put]
func UpdateCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		categoryID := c.Param("id")
		var category model.Category

		if err := db.Where("id = ? AND account_id = ?", categoryID, accountID).First(&category).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Category not found"})
			return
		}

		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request data"})
			return
		}

		category.AccountID = accountID.(uint)
		if err := db.Save(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to update category"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Category updated successfully"})
	}
}

// GetCategories godoc
// @Summary Get all categories
// @Description Retrieve all categories or filter by query parameters
// @Tags categories
// @Produce json
// @Success 200 {object} model.CategoriesResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /categories [get]
func GetCategories(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		var categories []model.Category
		if err := db.Where("account_id = ?", accountID).Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to retrieve categories"})
			return
		}

		var categoryResponses []model.CategoryResponse
		for _, category := range categories {
			categoryResponses = append(categoryResponses, model.CategoryResponse{
				ID:       category.ID,
				Name:     category.Name,
				ParentID: category.ParentID,
			})
		}

		c.JSON(http.StatusOK, model.CategoriesResponse{
			Message:    "Categories retrieved successfully",
			Categories: categoryResponses,
		})
	}
}

// SoftDeleteCategory godoc
// @Summary Soft delete a category
// @Description Soft deletes a category and reassigns its products to the default category
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /categories/{id} [delete]
func SoftDeleteCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		id := c.Param("id")
		var category model.Category
		if err := db.Where("id = ? AND account_id = ?", id, accountID).First(&category).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Category not found"})
			return
		}

		var defaultCategory model.Category
		if err := db.First(&defaultCategory, "name = ? AND account_id = ?", "Uncategorized", accountID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Default category not found"})
			return
		}

		if err := db.Model(&model.Product{}).Where("category_id = ? AND account_id = ?", category.ID, accountID).Update("category_id", defaultCategory.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to reassign products"})
			return
		}

		if err := db.Delete(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to delete category"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Category deleted successfully"})
	}
}

// HardDeleteCategory godoc
// @Summary Hard delete a category
// @Description Hard deletes a category and reassigns its products to the default category
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /categories/{id}/hard [delete]
func HardDeleteCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		id := c.Param("id")
		var category model.Category
		if err := db.Unscoped().Where("id = ? AND account_id = ?", id, accountID).First(&category).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Category not found"})
			return
		}

		var defaultCategory model.Category
		if err := db.First(&defaultCategory, "name = ? AND account_id = ?", "Uncategorized", accountID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Default category not found"})
			return
		}

		if err := db.Model(&model.Product{}).Where("category_id = ? AND account_id = ?", category.ID, accountID).Update("category_id", defaultCategory.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to reassign products"})
			return
		}

		if err := db.Unscoped().Delete(&category).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to delete category"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Category deleted successfully"})
	}
}

// RecoverCategory godoc
// @Summary Recover a deleted category
// @Description Recover a soft-deleted category by ID and reassign its products back to the category
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /categories/{id}/recover [post]
func RecoverCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		categoryID := c.Param("id")
		if result := db.Unscoped().Model(&model.Category{}).Where("id = ? AND account_id = ?", categoryID, accountID).Update("deleted_at", nil); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to recover category"})
			return
		}

		if err := db.Model(&model.Product{}).Where("category_id IS NULL AND account_id = ? AND previous_category_id = ?", accountID, categoryID).Update("category_id", categoryID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to reassign products"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Category recovered successfully"})
	}
}
