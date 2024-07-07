package main

import (
	"integration-service/internal/config"
	"integration-service/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.LoadConfig()

	r := gin.Default()

	err := r.SetTrustedProxies([]string{"127.0.0.1", "192.168.1.0/24"})
	if err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	r.POST("/orders", service.CreateOrder)
	r.GET("/orders/:id", service.GetOrder) // Assuming there is a GetOrder service
	r.POST("/reports/sales", service.GetSalesReports)
	r.POST("/reports/inventory", service.GetInventoryLevel)
	r.POST("/reports/shipping", service.GetShippingStatus)
	r.POST("/reports/users", service.GetUserActivities)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("integration service running on port 8085")
	if err := r.Run(":8085"); err != nil {
		log.Fatal("failed to run server: ", err)
	}
}
