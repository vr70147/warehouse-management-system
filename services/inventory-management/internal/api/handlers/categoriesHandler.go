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
		var category model.Category

		// Bind JSON body to category struct
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Failed to read body",
			})
			return
		}

		// Save the new category to the database
		if result := db.Create(&category); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to create category",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Category created successfully",
		})
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
		categoryID := c.Param("id")
		var category model.Category

		// Find the category by ID
		if result := db.First(&category, categoryID); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "Category not found",
			})
			return
		}

		// Bind JSON body to category struct
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Invalid request data",
			})
			return
		}

		// Save the updated category to the database
		if result := db.Save(&category); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to update category",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Category updated successfully",
		})
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
		var categories []model.Category

		// Retrieve all categories from the database
		if result := db.Find(&categories); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to retrieve categories",
			})
			return
		}

		// Map categories to response format
		var categoryResponses []model.CategoryResponse
		for _, category := range categories {
			categoryResponses = append(categoryResponses, model.CategoryResponse{
				ID:       category.ID,
				Name:     category.Name,
				ParentID: nil, // Adjust based on your category hierarchy implementation
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
		var category model.Category
		id := c.Param("id")

		// Find the category by ID
		if err := db.First(&category, id).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Category not found"})
			return
		}

		// Get the default category
		var defaultCategory model.Category
		if err := db.First(&defaultCategory, "name = ?", "Uncategorized").Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Default category not found"})
			return
		}

		// Reassign products to the default category
		if err := db.Model(&model.Product{}).Where("category_id = ?", category.ID).Update("category_id", defaultCategory.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to reassign products"})
			return
		}

		// Soft delete the category
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
		var category model.Category
		id := c.Param("id")

		// Find the category by ID
		if err := db.Unscoped().First(&category, id).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Category not found"})
			return
		}

		// Get the default category
		var defaultCategory model.Category
		if err := db.First(&defaultCategory, "name = ?", "Uncategorized").Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Default category not found"})
			return
		}

		// Reassign products to the default category
		if err := db.Model(&model.Product{}).Where("category_id = ?", category.ID).Update("category_id", defaultCategory.ID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to reassign products"})
			return
		}

		// Hard delete the category
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
		categoryID := c.Param("id")

		// Recover the soft-deleted category by setting deleted_at to NULL
		if result := db.Model(&model.Category{}).Unscoped().Where("id = ?", categoryID).Update("deleted_at", nil); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to recover category",
			})
			return
		}

		// Reassign products back to the recovered category
		if err := db.Model(&model.Product{}).Where("category_id IS NULL").Where("previous_category_id = ?", categoryID).Update("category_id", categoryID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to reassign products"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Category recovered successfully",
		})
	}
}
