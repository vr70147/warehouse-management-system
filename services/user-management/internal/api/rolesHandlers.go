package api

import (
	"encoding/json"
	"net/http"
	"user-management/internal/initializers"
	"user-management/internal/model"

	"github.com/gin-gonic/gin"
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
