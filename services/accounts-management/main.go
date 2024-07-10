package main

import (
	"accounts-management/internal/api/routes"
	"accounts-management/internal/cache"
	"accounts-management/internal/initializers"
	"accounts-management/internal/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	cache.InitRedis()
}

func main() {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	routes.Routers(r, initializers.DB)

	r.Run()
}
