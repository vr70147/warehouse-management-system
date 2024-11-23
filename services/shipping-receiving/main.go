package main

import (
	_ "shipping-receiving/docs"
	"shipping-receiving/internal/api/routes"
	"shipping-receiving/internal/cache"
	"shipping-receiving/internal/initializers"
	"shipping-receiving/internal/kafka"
	"shipping-receiving/internal/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	cache.InitRedis()
}

func main() {
	go kafka.ConsumerShippingEvents()

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.Routers(r, initializers.DB, &utils.NotificationService{})

	r.Run()

}
