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
	router.GET("/users", api.RequireAuth, api.GetUsers)
	router.PUT("/users/:id", api.RequireAuth, api.UpdateUser)
	router.POST("/roles", api.RequireAuth, api.CreateRole)
	router.PUT("/roles/:id", api.RequireAuth, api.UpdateRole)
	router.GET("/roles", api.RequireAuth, api.GetRoles)
	router.DELETE("/roles/:id", api.RequireAuth, api.DeleteRole)
	router.PATCH("/roles/recover/:id", api.RequireAuth, api.RecoverRole)
	router.Run()
}
