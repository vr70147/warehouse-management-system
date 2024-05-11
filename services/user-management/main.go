package main

import (
	"user-management/internal/api"
	"user-management/internal/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabse()
}

func main() {

	router := gin.Default()

	//initial Postgres
	router.POST("/signup", api.Signup)
	router.POST("/login", api.Login)
	router.GET("/validate", api.RequireAuth, api.Validate)
	router.POST("/logout", api.Logout)
	router.Run()
}
