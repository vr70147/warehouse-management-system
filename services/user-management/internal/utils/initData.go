package utils

import (
	"user-management/internal/model"

	"gorm.io/gorm"
)

// Initial data for roles and departments
func InitDefaultData(db *gorm.DB) {
	var count int64
	db.Model(&model.Department{}).Count(&count)
	if count == 0 {
		// Create default department
		department := model.Department{
			Name: "Default Department",
		}
		db.Create(&department)

		// Create default role
		role := model.Role{
			Role:         "super-admin",
			Description:  "Super Admin",
			DepartmentID: department.ID,
		}
		db.Create(&role)
	}
}
