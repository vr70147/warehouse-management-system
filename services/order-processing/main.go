package main

import (
	_ "order-processing/docs"
	"order-processing/internal/api/routes"
	"order-processing/internal/cache"
	"order-processing/internal/initializers"
	"order-processing/internal/kafka"
	"order-processing/internal/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Initialize environment variables, database connection, Redis cache, and Kafka writers.
func init() {
	initializers.LoadEnvVariables() // Load environment variables from .env file
	initializers.ConnectToDB()      // Establish database connection
	cache.InitRedis()               // Initialize Redis cache

	kafka.InitKafkaWriters() // Initialize Kafka writers for various topics
}

func main() {
	// Start Kafka consumers in separate goroutines
	go kafka.ConsumerOrderEvent()      // Consume order events
	go kafka.ConsumerInventoryStatus() // Consume inventory status updates
	go kafka.ConsumerShippingStatus()  // Consume shipping status updates

	// Initialize Gin router
	r := gin.Default()

	// Setup Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup API routes
	routes.Routers(r, initializers.DB, &utils.NotificationService{})

	// Run the Gin server
	r.Run() // Default listens on :8083
}
