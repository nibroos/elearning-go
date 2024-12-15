package config

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/dtos"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{Client: client}
}

func FetchCachedData(ctx context.Context, sqlDB *sqlx.DB) error {
	err := FetchAndCacheSubscribes(ctx, sqlDB)
	if err != nil {
		return err
	}

	return nil
}

func FetchAndCacheSubscribes(ctx context.Context, sqlDB *sqlx.DB) error {
	var subscribes []dtos.SubscribeListDTO

	query := `SELECT s.id, s.name, s.description, s.created_at, s.updated_at, s.deleted_at,
        cu.name as created_by_name,
        uu.name as updated_by_name
    FROM subscribes s
    JOIN users cu ON s.created_by_id = cu.id
    LEFT JOIN users uu ON s.updated_by_id = uu.id
    WHERE s.deleted_at IS NULL`

	err := sqlDB.SelectContext(ctx, &subscribes, query)
	if err != nil {
		return err
	}

	// Marshal the data to JSON
	data, err := json.Marshal(subscribes)
	if err != nil {
		return err
	}

	// Store the data in Redis
	err = RedisClient.Set(ctx, "subscribes", data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// func FetchAndCacheSubscribes(ctx context.Context, gormDB *gorm.DB, sqlDB *sqlx.DB, redisCache *cache.RedisCache) error {
// 	repo := repository.NewSubscribeRepository(gormDB, sqlDB, redisCache)
// 	err := repo.FetchAndCacheSubscribes(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
