package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl string
	JwtSecret   string
	Port        string
}

var AppConfig Config

func LoadConfig(configPath ...string) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	AppConfig = Config{
		DatabaseUrl: getEnv("DATABASE_URL", ""),
		JwtSecret:   getEnv("JWT_SECRET", ""),
		Port:        getEnv("PORT", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return defaultValue
}
