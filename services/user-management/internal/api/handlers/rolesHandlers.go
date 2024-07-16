package handlers

import (
	"net/http"
	"user-management/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role with permissions and department
// @Tags roles
// @Accept json
// @Produce json
// @Param body body model.Role true "Role data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /roles [post]
func CreateRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		var role model.Role

		if err := c.ShouldBindJSON(&role); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Failed to read body"})
			return
		}

		role.AccountID = accountID.(uint)

		if result := db.Create(&role); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to create role"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Role created successfully", Data: role})
	}
}

// UpdateRole godoc
// @Summary Update a role
// @Description Update a role by ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param body body model.Role true "Role data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /roles/{id} [put]
func UpdateRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		roleID := c.Param("id")

		var body model.Role

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request data"})
			return
		}

		var role model.Role
		if result := db.Where("id = ? AND account_id = ?", roleID, accountID).First(&role); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Role not found"})
			return
		}

		if body.Role != "" {
			role.Role = body.Role
		}
		if body.Description != "" {
			role.Description = body.Description
		}
		if body.IsActive {
			role.IsActive = body.IsActive
		}
		if body.DepartmentID != 0 {
			role.DepartmentID = body.DepartmentID
		}

		if result := db.Save(&role); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to update role"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Role updated successfully", Data: role})
	}
}

// GetRoles godoc
// @Summary Get all roles
// @Description Retrieve all roles or filter by query parameters
// @Tags roles
// @Produce json
// @Param role_name query string false "Role name"
// @Success 200 {object} model.RolesResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /roles [get]
func GetRoles(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		queryCondition := model.Role{AccountID: accountID.(uint)}
		var queryExists bool

		if role := c.Query("role_name"); role != "" {
			queryCondition.Role = role
			queryExists = true
		}

		var roles []model.Role
		var result *gorm.DB

		if queryExists {
			result = db.Where(&queryCondition).Find(&roles)
		} else {
			result = db.Where("account_id = ?", accountID).Find(&roles)
		}

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to retrieve roles"})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "No roles found matching the criteria"})
			return
		}

		c.JSON(http.StatusOK, model.RolesResponse{Message: "Roles retrieved successfully", Roles: roles})
	}
}

// SoftDeleteRole godoc
// @Summary Soft delete a role
// @Description Soft delete a role by ID
// @Tags roles
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /roles/{id} [delete]
func SoftDeleteRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		roleID := c.Param("id")

		var role model.Role
		if result := db.Where("id = ? AND account_id = ?", roleID, accountID).First(&role); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Role not found"})
			return
		}

		if result := db.Delete(&role); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to soft delete role"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Role soft deleted successfully"})
	}
}

// HardDeleteRole godoc
// @Summary Hard delete a role
// @Description Hard delete a role by ID
// @Tags roles
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /roles/hard/{id} [delete]
func HardDeleteRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		roleID := c.Param("id")

		if result := db.Unscoped().Where("id = ? AND account_id = ?", roleID, accountID).Delete(&model.Role{}); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to hard delete role"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Role hard deleted successfully"})
	}
}

// RecoverRole godoc
// @Summary Recover a deleted role
// @Description Recover a soft-deleted role by ID
// @Tags roles
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /roles/{id}/recover [post]
func RecoverRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		roleID := c.Param("id")

		var role model.Role
		if result := db.Unscoped().Where("id = ? AND account_id = ?", roleID, accountID).First(&role); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Role not found"})
			return
		}

		if result := db.Model(&role).Update("deleted_at", nil); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to recover role"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Role recovered successfully"})
	}
}
