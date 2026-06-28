package controller

import (
	"context"
	"time"
)

type mockRepository struct {
	redis *MockRedisClient
}

func (m *mockRepository) List(ctx context.Context, key string) ([]string, error) {
	return m.redis.LRange(ctx, key, 0, -1).Result()
}

func (m *mockRepository) Save(ctx context.Context, key string, value string) error {
	_, err := m.redis.LPush(ctx, key, value).Result()
	return err
}

func (m *mockRepository) Delete(ctx context.Context, keys ...string) error {
	_, err := m.redis.Del(ctx, keys...).Result()
	return err
}

func (m *mockRepository) Get(ctx context.Context, key string) (string, error) {
	return m.redis.Get(ctx, key).Result()
}

func (m *mockRepository) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return m.redis.Set(ctx, key, value, expiration).Err()
}
