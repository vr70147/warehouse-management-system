package main

import (
	_ "order-processing/docs"
	"order-processing/internal/api/routes"
	"order-processing/internal/cache"
	"order-processing/internal/initializers"
	"order-processing/kafka"

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
	go kafka.ConsumerOrderEvent()

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.Routers(r, initializers.DB)

	r.Run()

}
