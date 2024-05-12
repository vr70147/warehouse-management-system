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

func CreateRole(c *gin.Context) {
	var body struct {
		RoleName    string `gorm:"unique:not null"`
		Description string
		Permission  map[string]interface{} `json:"permission" binding:"required"`
		IsActive    bool                   `gorm:"default:true"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to read body",
		})
		return
	}
	permissionsJSON, err := json.Marshal(body.Permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to encode permissions",
		})
		return
	}
	role := model.Roles{
		RoleName:    body.RoleName,
		Description: body.Description,
		Permission:  string(permissionsJSON),
		IsActive:    body.IsActive,
	}
	result := initializers.DB.Create(&role)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create role",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func UpdateRole(c *gin.Context) {
	roleID := c.Param("id")
	log.Printf("Received roleID: %s", roleID)

	var body struct {
		RoleName    string                 `json:"roleName"`
		Description string                 `json:"description"`
		Permission  map[string]interface{} `json:"permission"`
		IsActive    *bool                  `json:"isActive"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}
	updateData := make(map[string]interface{})
	if body.RoleName != "" {
		updateData["role_name"] = body.RoleName
	}
	if body.Description != "" {
		updateData["description"] = body.Description
	}
	if body.Permission != nil {
		permissionJSON, err := json.Marshal(body.Permission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to encode permissions",
			})
			return
		}
		updateData["permission"] = gorm.Expr("permission || ?", string(permissionJSON))
	}
	if body.IsActive != nil {
		updateData["is_active"] = body.IsActive
	}

	result := initializers.DB.Model(&model.Roles{}).Where("id = ?", roleID).Updates(updateData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update role",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})

}

func GetRoles(c *gin.Context) {
	queryCondition := model.Roles{}

	var queryExists bool

	if roleName := c.Query("role_name"); roleName != "" {
		queryCondition.RoleName = roleName
		queryExists = true
	}

	var roles []model.Roles
	var result *gorm.DB

	if queryExists || len(c.Params) == 0 {
		result = initializers.DB.Where(&queryCondition).Find(&roles)
	} else {
		result = initializers.DB.Find(&roles)
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve roles",
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No roles found matching the criteria",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
	})
}

func DeleteRole(c *gin.Context) {
	roleID := c.Param("id")

	result := initializers.DB.Delete(&model.Roles{}, roleID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete role",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No role found with the given ID",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfuly"})
}

func RecoverRole(c *gin.Context) {
	roleID := c.Param("id")

	result := initializers.DB.Model(&model.Roles{}).Unscoped().Where("id = ?", roleID).Update("deleted_at", nil)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to recover role",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No deleted role found with the given ID",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role recovered successfully"})
}
