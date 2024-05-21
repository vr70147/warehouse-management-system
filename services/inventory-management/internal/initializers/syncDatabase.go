package initializers

import (
	"inventory-management/internal/model"
)

func SyncDatabse() {
	DB.AutoMigrate(&model.Product{}, &model.Stock{}, &model.Category{}, &model.Supplier{})
}
