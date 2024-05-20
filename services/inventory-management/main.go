package main

import (
	"inventory-management/internal/api"
	"inventory-management/internal/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabse()
}

func main() {
	api.Routers()

}
