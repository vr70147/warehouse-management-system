package main

import (
	"log"
	_ "user-management/docs"
	"user-management/internal/api/routes"
	"user-management/internal/cache"
	"user-management/internal/initializers"
	"user-management/internal/utils"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabse()
	cache.InitRedis()
	utils.InitDefaultData(initializers.DB)

	// Initialize the default email sender
	emailSender := &utils.DefaultEmailSender{}
	scheduler := utils.NewScheduler(initializers.DB, emailSender)
	scheduler.StartMonthlySummaryScheduler()
}

// @title User Management API
// @version 1.0
// @description This is a user management server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()

	if err := r.SetTrustedProxies([]string{"127.0.0.1", "192.168.1.1"}); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.Routers(r, initializers.DB)

	r.Run()
}
