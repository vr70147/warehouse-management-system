package main

import (
	"accounts-management/internal/cache"
	"accounts-management/internal/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	cache.InitRedis()
}
