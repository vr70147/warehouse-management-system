package main

import (
	"context"
	"log"
	"os"
	"user-management/internal/api"
	"user-management/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func main() {

	config.LoadConfig()
	dbUrl := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	router := gin.Default()

	router.POST("/user", api.Signup(conn))
	if err := router.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
