package initializers

import (
	"inventory-management/internal/model"
)

func SyncDatabse() {
	DB.AutoMigrate(&model.InventoryItem{})
	DB.AutoMigrate(&model.Location{})
	DB.AutoMigrate(&model.StockMovement{})
	DB.AutoMigrate(&model.Supplier{})
}
