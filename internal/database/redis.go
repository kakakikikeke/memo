package database

import (
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Close() error
}

func NewClient() *redis.Client {
	redisURL := "redis://localhost:6379/0"
	if url := os.Getenv("REDIS_URL"); url != "" {
		redisURL = url
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(opt)
}
