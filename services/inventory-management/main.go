package main

import (
	"fmt"
	_ "inventory-management/docs"
	"inventory-management/internal/api/routes"
	"inventory-management/internal/cache"
	"inventory-management/internal/initializers"
	"inventory-management/internal/kafka"
	"inventory-management/internal/utils"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	cache.InitRedis()
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
	fmt.Println("Environment Variables:")
	for _, e := range os.Environ() {
		fmt.Println(e)
	}

	go kafka.ConsumerOrderEvents()
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.Routers(r, initializers.DB, &utils.NotificationService{})

	r.Run()

}
