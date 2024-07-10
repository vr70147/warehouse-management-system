package api

import (
	"encoding/json"
	"log"
	"net/http"
	"user-management/internal/initializers"
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
		var body struct {
			Role         string `gorm:"unique;not null"`
			Description  string
			Permissions  []model.Permission `gorm:"many2many:role_permissions"`
			IsActive     bool               `gorm:"default:true"`
			Users        []model.User       `gorm:"foreignKey:RoleID"`
			DepartmentID uint               `gorm:"not null"`
			Department   model.Department   `gorm:"foreignKey:DepartmentID"`
		}

		if c.Bind(&body) != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Failed to read body",
			})
			return
		}

		role := model.Role{
			Role:         body.Role,
			Description:  body.Description,
			Permissions:  body.Permissions,
			IsActive:     body.IsActive,
			DepartmentID: body.DepartmentID,
		}
		result := initializers.DB.Create(&role)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Failed to create role",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Role created successfully",
		})
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
		roleID := c.Param("id")
		log.Printf("Received roleID: %s", roleID)

		var body struct {
			Role         string `gorm:"unique;not null"`
			Description  string
			Permission   map[string]interface{} `json:"permission"`
			IsActive     bool                   `gorm:"default:true"`
			Users        []model.User           `gorm:"foreignKey:RoleID"`
			DepartmentID uint                   `gorm:"not null"`
			Department   model.Department       `gorm:"foreignKey:DepartmentID"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Invalid request data",
			})
			return
		}
		updateData := make(map[string]interface{})
		if body.Role != "" {
			updateData["role"] = body.Role
		}
		if body.Description != "" {
			updateData["description"] = body.Description
		}
		if body.Permission != nil {
			permissionJSON, err := json.Marshal(body.Permission)
			if err != nil {
				c.JSON(http.StatusInternalServerError, model.ErrorResponse{
					Error: "Failed to encode permissions",
				})
				return
			}
			updateData["permission"] = gorm.Expr("permission || ?", string(permissionJSON))
		}

		result := initializers.DB.Model(&model.Role{}).Where("id = ?", roleID).Updates(updateData)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to update role",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Role updated successfully",
		})
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
		queryCondition := model.Role{}

		var queryExists bool

		if role := c.Query("role_name"); role != "" {
			queryCondition.Role = role
			queryExists = true
		}

		var roles []model.Role
		var result *gorm.DB

		if queryExists || len(c.Params) == 0 {
			result = initializers.DB.Where(&queryCondition).Find(&roles)
		} else {
			result = initializers.DB.Find(&roles)
		}

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to retrieve roles",
			})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "No roles found matching the criteria",
			})
			return
		}

		c.JSON(http.StatusOK, model.RolesResponse{
			Message: "Roles retrieved successfully",
			Roles:   roles,
		})
	}
}

// DeleteRole godoc
// @Summary Delete a role
// @Description Delete a role by ID
// @Tags roles
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /roles/{id} [delete]
func DeleteRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := c.Param("id")

		result := initializers.DB.Delete(&model.Role{}, roleID)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to delete role",
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "No role found with the given ID",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Role deleted successfully",
		})
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
		roleID := c.Param("id")

		result := initializers.DB.Model(&model.Role{}).Unscoped().Where("id = ?", roleID).Update("deleted_at", nil)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to recover role",
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "No deleted role found with the given ID",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Role recovered successfully",
		})
	}
}
