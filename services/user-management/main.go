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
	router.Run()

	// router.POST("/user", api.Signup(conn))
	// if err := router.Run(":" + os.Getenv("PORT")); err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }

}
