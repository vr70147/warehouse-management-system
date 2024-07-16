package main

import (
	"accounts-management/internal/api/routes"
	"accounts-management/internal/cache"
	"accounts-management/internal/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	cache.InitRedis()
}

func main() {
	r := gin.Default()
	routes.Routers(r, initializers.DB)

	r.Run()
}
