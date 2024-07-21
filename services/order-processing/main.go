package main

import (
	_ "order-processing/docs"
	"order-processing/internal/api/routes"
	"order-processing/internal/cache"
	"order-processing/internal/initializers"
	"order-processing/internal/utils"
	message_broker "order-processing/message-broker"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	cache.InitRedis()

	message_broker.KafkaWriterInstance = message_broker.NewKafkaWriter(
		[]string{os.Getenv("KAFKA_BROKERS")},
		os.Getenv("ORDER_EVENTS_TOPIC"),
	)
}

func main() {

	go message_broker.ConsumerOrderEvent()

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.Routers(r, initializers.DB, &utils.NotificationService{})

	r.Run()

}
