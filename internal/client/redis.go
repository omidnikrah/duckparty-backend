package client

import (
	"github.com/omidnikrah/duckparty-backend/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(config *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword,
	})
}
