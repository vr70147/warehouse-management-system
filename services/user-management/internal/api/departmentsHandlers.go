package api

import (
	"net/http"
	"user-management/internal/initializers"
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
		var body struct {
			Name  string
			Roles []model.Role
		}

		if c.Bind(&body) != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Failed to read body",
			})
			return
		}
		department := model.Department{
			Name: body.Name,
		}

		result := initializers.DB.Create(&department)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Failed to create department",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Department created successfully",
		})
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
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /departments/{id} [put]
func UpdateDepartment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		departmentID := c.Param("id")

		var body struct {
			Name  string
			Roles []model.Role
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Invalid request data",
			})
			return
		}
		updateData := make(map[string]interface{})
		if body.Name != "" {
			updateData["name"] = body.Name
		}

		result := initializers.DB.Model(&model.Department{}).Where("id = ?", departmentID).Updates(updateData)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to update department",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Department updated successfully",
		})
	}
}

// GetDepartment godoc
// @Summary Get all departments
// @Description Retrieve all departments or filter by query parameters
// @Tags departments
// @Produce json
// @Param name query string false "Name"
// @Success 200 {object} model.DepartmentsResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /departments [get]
func GetDepartment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryCondition := model.Department{}
		var queryExists bool

		if name := c.Query("name"); name != "" {
			queryCondition.Name = name
			queryExists = true
		}

		var departments []model.Department
		var result *gorm.DB

		if queryExists {
			result = initializers.DB.Where(&queryCondition).Find(&departments)
		} else {
			result = initializers.DB.Find(&departments)
		}

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to retrieve departments",
			})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "No departments found matching the criteria",
			})
			return
		}

		// Fetch roles for each department
		var departmentResponses []model.DepartmentResponse
		for _, department := range departments {
			var roles []model.Role
			initializers.DB.Model(&model.Role{}).Where("department_id = ?", department.ID).Find(&roles)
			departmentResponses = append(departmentResponses, model.DepartmentResponse{
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

// DeleteDepartment godoc
// @Summary Delete a department
// @Description Delete a department by ID
// @Tags departments
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /departments/{id} [delete]
func DeleteDepartment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		departmentID := c.Param("id")

		result := initializers.DB.Delete(&model.Department{}, departmentID)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to delete department",
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "No department found with the given ID",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Department deleted successfully",
		})
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
		departmentID := c.Param("id")

		result := initializers.DB.Model(&model.Department{}).Unscoped().Where("id = ?", departmentID).Update("deleted_at", nil)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to recover department",
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "No deleted department found with the given ID",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Department recovered successfully",
		})
	}
}

// GetUsersByDepartment godoc
// @Summary Get users by department
// @Description Retrieve users by department name
// @Tags departments
// @Produce json
// @Param department query string true "Department name"
// @Success 200 {object} model.UsersResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /departments/users [get]
func GetUsersByDepartment(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		departmentName := c.Query("department")
		if departmentName == "" {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Department name is required",
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
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to retrieve users: " + result.Error.Error(),
			})
			return
		}

		if len(users) == 0 {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "No users found in the specified department",
			})
			return
		}

		var userResponses []model.UserResponse
		for _, user := range users {
			userResponses = append(userResponses, model.UserResponse{
				User:       user.User,
				RoleID:     user.Role,
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
