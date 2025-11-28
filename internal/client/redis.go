package client

import (
	"fmt"

	"github.com/omidnikrah/duckparty-backend/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(config *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
	})
}
