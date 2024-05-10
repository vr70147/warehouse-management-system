package config

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

var appConfig Config

func LoadConfig(configPaths ...string) {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	appConfig = Config{
		DatabaseUrl: os.Getenv("DATABASE_URL"),
		JwtSecret:   os.Getenv("JWT_SECRET"),
		Port:        os.Getenv("PORT"),
	}

}
