package main

import (
	"os"
	"user-management/internal/config"
)

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	config.ConnectToDB(dbUrl)
}
