package api

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

		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Failed to read body",
			})
			return
		}

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

		if result := db.First(&category, categoryID); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "Category not found",
			})
			return
		}

		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Invalid request data",
			})
			return
		}

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
		if result := db.Find(&categories); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to retrieve categories",
			})
			return
		}

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

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /categories/{id} [delete]
func DeleteCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Param("id")

		if result := db.Delete(&model.Category{}, categoryID); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to delete category",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Category deleted successfully",
		})
	}
}

// RecoverCategory godoc
// @Summary Recover a deleted category
// @Description Recover a soft-deleted category by ID
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

		if result := db.Model(&model.Category{}).Unscoped().Where("id = ?", categoryID).Update("deleted_at", nil); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to recover category",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Category recovered successfully",
		})
	}
}

// HardDeleteCategory godoc
// @Summary Hard-delete a category
// @Description Hard-delete a category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 500 {object} model.ErrorResponse
func HardDeleteCategory(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Param("id")

		if result := db.Unscoped().Delete(&model.Category{}, categoryID); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to delete category",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Category deleted permanently",
		})
	}
}
