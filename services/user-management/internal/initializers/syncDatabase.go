package initializers

import (
	"user-management/internal/model"
)

func SyncDatabse() {
	DB.AutoMigrate(&model.User{})
	DB.AutoMigrate(&model.Roles{})
}
