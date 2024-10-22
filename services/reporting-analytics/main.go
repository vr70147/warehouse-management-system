package main

import (
	_ "reporting-analytics/docs"
	"reporting-analytics/internal/api/routes"
	"reporting-analytics/internal/cache"
	"reporting-analytics/internal/initializers"
	"reporting-analytics/internal/kafka"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	cache.InitRedis()
	kafka.InitKafkaWriters()
}

// @title Reporting Analytics Service API
// @version 1.0
// @description This is a server for reporting analytics service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:8084
// @BasePath /
func main() {
	go kafka.ConsumerSalesEvent()
	go kafka.ConsumerInventoryLevel()
	go kafka.ConsumerShippingStatus()
	go kafka.ConsumerUserActivity()

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.Routers(r, initializers.DB)

	r.Run()

}
