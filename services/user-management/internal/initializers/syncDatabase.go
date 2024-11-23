package initializers

import (
	"user-management/internal/model"
)

func SyncDatabse() {
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Role{})
	DB.AutoMigrate(&model.Department{})
	DB.AutoMigrate(&model.Role{}, &model.User{})
}
