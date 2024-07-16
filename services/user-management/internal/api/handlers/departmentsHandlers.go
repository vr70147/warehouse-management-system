package handlers

import (
	"net/http"
	"user-management/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateDepartment godoc
// @Summary Create a new department
// @Description Create a new department
// @Tags departments
// @Accept json
// @Produce json
// @Param body body model.Department true "Department data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /departments [post]
func CreateDepartment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		var body struct {
			Name  string
			Roles []model.Role
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request data"})
			return
		}

		department := model.Department{
			Name:      body.Name,
			AccountID: accountID.(uint),
		}

		if result := db.Create(&department); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to create department"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Department created successfully", Data: department})
	}
}

// UpdateDepartment godoc
// @Summary Update a department
// @Description Update a department by ID
// @Tags departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param body body model.Department true "Department data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /departments/{id} [put]
func UpdateDepartment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		departmentID := c.Param("id")

		var body struct {
			Name  string
			Roles []model.Role
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request data"})
			return
		}

		var department model.Department
		if result := db.Where("id = ? AND account_id = ?", departmentID, accountID).First(&department); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Department not found"})
			return
		}

		if body.Name != "" {
			department.Name = body.Name
		}

		if result := db.Save(&department); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to update department"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Department updated successfully", Data: department})
	}
}

// GetDepartments godoc
// @Summary Get all departments
// @Description Retrieve all departments or filter by query parameters
// @Tags departments
// @Produce json
// @Param name query string false "Name"
// @Success 200 {object} model.DepartmentsResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /departments [get]
func GetDepartments(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		queryCondition := model.Department{AccountID: accountID.(uint)}
		var queryExists bool

		if name := c.Query("name"); name != "" {
			queryCondition.Name = name
			queryExists = true
		}

		var departments []model.Department
		var result *gorm.DB

		if queryExists {
			result = db.Where(&queryCondition).Find(&departments)
		} else {
			result = db.Where("account_id = ?", accountID).Find(&departments)
		}

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to retrieve departments"})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "No departments found matching the criteria"})
			return
		}

		// Fetch roles for each department
		var departmentResponses []model.DepartmentResponse
		for _, department := range departments {
			var roles []model.Role
			db.Model(&model.Role{}).Where("department_id = ?", department.ID).Find(&roles)
			departmentResponses = append(departmentResponses, model.DepartmentResponse{
				ID:    department.ID,
				Name:  department.Name,
				Roles: roles,
			})
		}

		c.JSON(http.StatusOK, model.DepartmentsResponse{
			Message:     "Departments retrieved successfully",
			Departments: departmentResponses,
		})
	}
}

// SoftDeleteDepartment godoc
// @Summary Soft delete a department
// @Description Soft delete a department by ID
// @Tags departments
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /departments/{id} [delete]
func SoftDeleteDepartment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		departmentID := c.Param("id")

		var department model.Department
		if result := db.Where("id = ? AND account_id = ?", departmentID, accountID).First(&department); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Department not found"})
			return
		}

		if result := db.Delete(&department); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to soft delete department"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Department soft deleted successfully"})
	}
}

// HardDeleteDepartment godoc
// @Summary Hard delete a department
// @Description Hard delete a department by ID
// @Tags departments
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /departments/hard/{id} [delete]
func HardDeleteDepartment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		departmentID := c.Param("id")

		if result := db.Unscoped().Where("id = ? AND account_id = ?", departmentID, accountID).Delete(&model.Department{}); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to hard delete department"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Department hard deleted successfully"})
	}
}

// RecoverDepartment godoc
// @Summary Recover a deleted department
// @Description Recover a soft-deleted department by ID
// @Tags departments
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /departments/{id}/recover [post]
func RecoverDepartment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		departmentID := c.Param("id")

		var department model.Department
		if result := db.Unscoped().Where("id = ? AND account_id = ?", departmentID, accountID).First(&department); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Department not found"})
			return
		}

		if result := db.Model(&department).Update("deleted_at", nil); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to recover department"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Department recovered successfully"})
	}
}

// GetUsersByDepartment godoc
// @Summary Get users by department
// @Description Retrieve users by department name
// @Tags departments
// @Produce json
// @Param department query string true "Department name"
// @Success 200 {object} model.UsersResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /departments/users [get]
func GetUsersByDepartment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		departmentName := c.Query("department")
		if departmentName == "" {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Department name is required"})
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
		result := db.Table("users").
			Select("users.*, roles.role as role, roles.permission, roles.is_active, departments.name as department_name").
			Joins("left join roles on roles.id = users.role_id").
			Joins("left join departments on departments.id = roles.department_id").
			Where("departments.name = ? AND departments.account_id = ?", departmentName, accountID).
			Scan(&users)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to retrieve users: " + result.Error.Error()})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "No users found in the specified department"})
			return
		}

		var userResponses []model.UserResponse
		for _, user := range users {
			userResponses = append(userResponses, model.UserResponse{
				User:       user.User,
				Role:       user.Role,
				Permission: user.Permission,
				IsActive:   user.IsActive,
				Department: user.DepartmentName,
			})
		}

		c.JSON(http.StatusOK, model.UsersResponse{
			Message: "Users retrieved successfully",
			Users:   userResponses,
		})
	}
}
