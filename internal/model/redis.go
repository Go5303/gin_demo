package model

import (
	"context"
	"fmt"
	"log"
	"time"

	"git.woda.ink/Woda_OA/config"
	"git.woda.ink/Woda_OA/pkg/cache"
	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var redisPrefix string

// InitRedis initializes the Redis connection
func InitRedis(cfg *config.RedisConfig) error {
	redisPrefix = cfg.Prefix

	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	// Initialize cache package
	cache.Init(RDB, redisPrefix)

	log.Println("[Redis] Connected successfully")
	return nil
}

// GetRedis returns the Redis client
func GetRedis() *redis.Client {
	return RDB
}

// RedisKey builds a prefixed cache key
func RedisKey(key string) string {
	if redisPrefix != "" {
		return redisPrefix + ":" + key
	}
	return key
}

// CloseRedis closes the Redis connection
func CloseRedis() error {
	if RDB != nil {
		return RDB.Close()
	}
	return nil
}

func GetKey(key string) string {
	return RDB.Get(context.Background(), key).Val()
}
