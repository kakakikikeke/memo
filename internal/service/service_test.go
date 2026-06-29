package service

import (
	"context"
	"testing"
	"time"

	"github.com/kakakikikeke/memo/internal/repository"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type inMemoryRepository struct {
	data  map[string][]string
	users map[string]string
}

func newInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{data: make(map[string][]string), users: make(map[string]string)}
}

func (r *inMemoryRepository) List(ctx context.Context, key string) ([]string, error) {
	return append([]string(nil), r.data[key]...), nil
}

func (r *inMemoryRepository) Save(ctx context.Context, key string, value string) error {
	r.data[key] = append([]string{value}, r.data[key]...)
	return nil
}

func (r *inMemoryRepository) Delete(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		delete(r.data, key)
	}
	return nil
}

func (r *inMemoryRepository) Get(ctx context.Context, key string) (string, error) {
	value, ok := r.users[key]
	if !ok {
		return "", redis.Nil
	}
	return value, nil
}

func (r *inMemoryRepository) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	r.users[key] = value
	return nil
}

func (r *inMemoryRepository) Text() repository.TextRepository {
	return r
}

func (r *inMemoryRepository) Image() repository.ImageRepository {
	return r
}

func (r *inMemoryRepository) File() repository.FileRepository {
	return r
}

func (r *inMemoryRepository) User() repository.UserRepository {
	return r
}

func TestSaveTextPersistsMessage(t *testing.T) {
	repo := newInMemoryRepository()
	svc := NewMemoService(repo)

	err := svc.SaveText(context.Background(), "alice", "hello")
	require.NoError(t, err)

	items, err := repo.List(context.Background(), "alice:memo")
	require.NoError(t, err)
	require.Equal(t, []string{"hello"}, items)
}

func TestCreateUserAndAuthenticateSuccess(t *testing.T) {
	repo := newInMemoryRepository()
	svc := NewMemoService(repo)

	err := svc.CreateUser(context.Background(), "alice", "secret")
	require.NoError(t, err)

	ok, err := svc.Authenticate(context.Background(), "alice", "secret")
	require.NoError(t, err)
	require.True(t, ok)
}

func TestAuthenticateWrongPassword(t *testing.T) {
	repo := newInMemoryRepository()
	svc := NewMemoService(repo)

	err := svc.CreateUser(context.Background(), "alice", "secret")
	require.NoError(t, err)

	ok, err := svc.Authenticate(context.Background(), "alice", "wrong")
	require.NoError(t, err)
	require.False(t, ok)
}
