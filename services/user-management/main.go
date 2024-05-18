package main

import (
	"user-management/internal/api"
	"user-management/internal/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabse()
}

func main() {
	api.Routers()
}
