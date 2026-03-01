package database

import (
	"context"
	"log"

	"my-portfolio/internal/config"

	"github.com/redis/go-redis/v9"
)

// InitRedis creates and verifies a Redis client from config.
// The server will not start if Redis is unreachable.
func InitRedis(cfg config.TypeMyPortfolio) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis connection failed (%s): %v", cfg.Redis.Addr, err)
	}

	log.Printf("Connected to Redis at %s (db=%d)", cfg.Redis.Addr, cfg.Redis.DB)
	return rdb
}
