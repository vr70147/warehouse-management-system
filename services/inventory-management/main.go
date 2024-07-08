package main

import (
	_ "inventory-management/docs"
	"inventory-management/internal/api/routes"
	"inventory-management/internal/initializers"
	"inventory-management/internal/middleware"
	"inventory-management/kafka"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

// @title Inventory API
// @version 1.0
// @description This is a server for managing inventory.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /
func main() {
	go kafka.ConsumerOrderEvents()
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.Routers(r, initializers.DB)

	r.Run()

}
