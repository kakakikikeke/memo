package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	List(ctx context.Context, key string) ([]string, error)
	Save(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, keys ...string) error
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}

type RedisClient interface {
	LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Close() error
}

type RedisRepository struct {
	client RedisClient
}

func NewRedisRepository(client RedisClient) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) List(ctx context.Context, key string) ([]string, error) {
	return r.client.LRange(ctx, key, 0, -1).Result()
}

func (r *RedisRepository) Save(ctx context.Context, key string, value string) error {
	_, err := r.client.LPush(ctx, key, value).Result()
	return err
}

func (r *RedisRepository) Delete(ctx context.Context, keys ...string) error {
	_, err := r.client.Del(ctx, keys...).Result()
	return err
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisRepository) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}
