package config

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedisClient() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Redis container name and port
		Password: "",           // No password set
		DB:       0,            // Use default DB
	})

	// Test the Redis connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}
