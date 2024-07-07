package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	OrderServiceURL       string
	OrderEventsTopic      string
	InventoryServiceURL   string
	ShippingServiceURL    string
	UserServiceURL        string
	ReportingAnalyticsURL string
	KafkaBrokers          string
	ReportingEventsTopic  string
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	OrderServiceURL = os.Getenv("ORDER_SERVICE_URL")
	OrderEventsTopic = os.Getenv("ORDER_EVENTS_TOPIC")
	InventoryServiceURL = os.Getenv("INVENTORY_SERVICE_URL")
	ShippingServiceURL = os.Getenv("SHIPPING_SERVICE_URL")
	UserServiceURL = os.Getenv("USER_SERVICE_URL")
	ReportingAnalyticsURL = os.Getenv("REPORTING_ANALYTICS_URL")
	KafkaBrokers = os.Getenv("KAFKA_BROKERS")
	ReportingEventsTopic = os.Getenv("REPORTING_EVENTS_TOPIC")
}
