package cache

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func SetCache(key string, value interface{}) error {
	err := redisClient.Set(ctx, key, value, 0).Err()
	return err
}

func GetCache(key string) (string, error) {
	val, err := redisClient.Get(ctx, key).Result()
	return val, err
}
