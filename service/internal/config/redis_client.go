package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedisClient() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: func() int {
			db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
			if err != nil {
				log.Fatalf("Invalid REDIS_DB value: %v", err)
			}
			return db
		}(),
	})

	// Test the Redis connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func InitRedisClientTest() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST_TEST"), os.Getenv("REDIS_PORT_TEST")),
		Password: os.Getenv("REDIS_PASSWORD_TEST"),
		DB: func() int {
			db, err := strconv.Atoi(os.Getenv("REDIS_DB_TEST"))
			if err != nil {
				log.Fatalf("Invalid REDIS_DB_TEST value: %v", err)
			}
			return db
		}(),
	})

	// Test the Redis connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}
