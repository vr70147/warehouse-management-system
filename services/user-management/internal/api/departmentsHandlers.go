package api

import (
	"net/http"
	"user-management/internal/initializers"
	"user-management/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateDepartment(c *gin.Context) {
	var body struct {
		Name  string
		Roles []model.Department
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to read body",
		})
		return
	}
	department := model.Department{
		Name: body.Name,
	}

	result := initializers.DB.Create(&department)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create department",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func UpdateDepartment(c *gin.Context) {
	departmentID := c.Param("id")

	var body struct {
		Name  string
		Roles []model.Department
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}
	updateData := make(map[string]interface{})
	if body.Name != "" {
		updateData["name"] = body.Name
	}

	result := initializers.DB.Model(&model.Department{}).Where("id = ?", departmentID).Updates(updateData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update department",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Department updated successfully"})
}

func GetDepartment(c *gin.Context) {
	queryCondition := model.Department{}

	var queryExists bool

	if name := c.Query("name"); name != "" {
		queryCondition.Name = name
		queryExists = true
	}

	var departments []model.Department
	var result *gorm.DB

	if queryExists || len(c.Params) == 0 {
		result = initializers.DB.Where(&queryCondition).Find(&departments)
	} else {
		result = initializers.DB.Find(&departments)
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve department",
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No department found matching the criteria",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"departments": departments,
	})
}

func DeleteDepartment(c *gin.Context) {
	departmentID := c.Param("id")

	result := initializers.DB.Delete(&model.Department{}, departmentID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete department",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No department found with the given ID",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Department deleted successfuly"})
}

func RecoverDepartment(c *gin.Context) {
	departmentID := c.Param("id")

	result := initializers.DB.Model(&model.Department{}).Unscoped().Where("id = ?", departmentID).Update("deleted_at", nil)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to recover department",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No deleted department found with the given ID",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "department recovered successfully"})
}

func GetUsersByDepartment(c *gin.Context) {
	departmentName := c.Query("department")
	if departmentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Department name is required",
		})
		return
	}

	var users []struct {
		model.User
		Role           string `gorm:"column:role"`
		Permission     string `gorm:"column:permission"`
		IsActive       bool   `gorm:"column:is_active"`
		DepartmentName string `gorm:"column:department_name"`
	}

	// Fetch users from the database
	result := initializers.DB.Table("users").
		Select("users.*, roles.role as role, roles.permission, roles.is_active, departments.name as department_name").
		Joins("left join roles on roles.id = users.role_id").
		Joins("left join departments on departments.id = roles.department_id").
		Where("departments.name = ?", departmentName).
		Scan(&users)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users: " + result.Error.Error(),
		})
		return
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No users found in the specified department",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
