package main

import (
	_ "customer-service/docs"
	"customer-service/internal/api/routes"
	"customer-service/internal/cache"
	"customer-service/internal/initializers"

	// "customer-service/internal/kafka"

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

	// go kafka.ConsumerOrderEvent()

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.Routers(r, initializers.DB)

	r.Run()

}
