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
		Role         string `gorm:"unique:not null"`
		Description  string
		Permission   map[string]interface{} `json:"permission"`
		IsActive     bool                   `gorm:"default:true"`
		Users        []model.User           `gorm:"foreignKey:RoleID"`
		DepartmentID uint                   `gorm:"not null"`
		Department   model.Department       `gorm:"foreignKey:DepartmentID"`
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
	role := model.Role{
		Role:         body.Role,
		Description:  body.Description,
		Permission:   string(permissionsJSON),
		IsActive:     body.IsActive,
		DepartmentID: body.DepartmentID,
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
		Role         string `gorm:"unique:not null"`
		Description  string
		Permission   map[string]interface{} `json:"permission"`
		IsActive     bool                   `gorm:"default:true"`
		Users        []model.User           `gorm:"foreignKey:RoleID"`
		DepartmentID uint                   `gorm:"not null"`
		Department   model.Department       `gorm:"foreignKey:DepartmentID"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to encode permissions",
			})
			return
		}
		updateData["permission"] = gorm.Expr("permission || ?", string(permissionJSON))
	}

	result := initializers.DB.Model(&model.Role{}).Where("id = ?", roleID).Updates(updateData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update role",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

func GetRoles(c *gin.Context) {
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

	result := initializers.DB.Delete(&model.Role{}, roleID)

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

	result := initializers.DB.Model(&model.Role{}).Unscoped().Where("id = ?", roleID).Update("deleted_at", nil)
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
