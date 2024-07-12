package main

import (
	"fmt"
	"integration-service/internal/config"
	"integration-service/internal/service"
	"log"
	"os"

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

	log.Printf("integration service running on port %s", os.Getenv("PORT"))
	if err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))); err != nil {
		log.Fatal("failed to run server: ", err)
	}
}
