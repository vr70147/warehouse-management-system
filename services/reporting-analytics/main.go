package main

import (
	_ "reporting-analytics/docs"
	"reporting-analytics/internal/api/routes"
	"reporting-analytics/internal/initializers"

	"reporting-analytics/internal/api/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.Routers(r, initializers.DB)

	r.Run()

}
