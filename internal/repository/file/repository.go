package file

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
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
