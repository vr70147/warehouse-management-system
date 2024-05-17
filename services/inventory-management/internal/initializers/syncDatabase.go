package initializers

import (
	"inventory-management/internal/model"
)

func SyncDatabse() {
	DB.AutoMigrate(&model.Item{})
	DB.AutoMigrate(&model.Order{})
	DB.AutoMigrate(&model.OrderItem{})
	DB.AutoMigrate(&model.Supplier{})
}
