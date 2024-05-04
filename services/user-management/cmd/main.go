package main

import (
	"log"
	"os"
	"user-management/configs"
	"user-management/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	configs.LoadConfig()
	startServer(configs.AppConfig.Port)
	dbUrl := os.Getenv("DATABASE_URL")
	config.ConnectToDB(dbUrl)
}

func startServer(port string) {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the User Management Service",
		})
	})
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
