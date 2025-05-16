package main

import (
	"time"

	"github.com/go-redis/redis"
)

type MockRedisClient struct {
	data map[string][]string
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{data: make(map[string][]string)}
}

func (m *MockRedisClient) LRange(key string, start, stop int64) *redis.StringSliceCmd {
	return redis.NewStringSliceResult(m.data[key], nil)
}

func (m *MockRedisClient) LPush(key string, values ...interface{}) *redis.IntCmd {
	for _, v := range values {
		m.data[key] = append([]string{v.(string)}, m.data[key]...)
	}
	return redis.NewIntResult(int64(len(m.data[key])), nil)
}

func (m *MockRedisClient) Del(keys ...string) *redis.IntCmd {
	for _, key := range keys {
		delete(m.data, key)
	}
	return redis.NewIntResult(1, nil)
}

func (m *MockRedisClient) Get(key string) *redis.StringCmd {
	return redis.NewStringResult("", redis.Nil)
}

func (m *MockRedisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return redis.NewStatusResult("OK", nil)
}

func (m *MockRedisClient) Close() error {
	return nil
}
