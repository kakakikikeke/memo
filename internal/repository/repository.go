package repository

import (
	"context"
	"time"

	filerepo "github.com/kakakikikeke/memo/internal/repository/file"
	imagerepo "github.com/kakakikikeke/memo/internal/repository/image"
	textrepo "github.com/kakakikikeke/memo/internal/repository/text"
	userrepo "github.com/kakakikikeke/memo/internal/repository/user"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	TextRepository
	ImageRepository
	FileRepository
	UserRepository
}

type Provider interface {
	Text() TextRepository
	Image() ImageRepository
	File() FileRepository
	User() UserRepository
}

type TextRepository interface {
	List(ctx context.Context, key string) ([]string, error)
	Save(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, keys ...string) error
}

type ImageRepository interface {
	List(ctx context.Context, key string) ([]string, error)
	Save(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, keys ...string) error
}

type FileRepository interface {
	List(ctx context.Context, key string) ([]string, error)
	Save(ctx context.Context, key string, value string) error
	Delete(ctx context.Context, keys ...string) error
}

type UserRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
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
	textRepo  TextRepository
	imageRepo ImageRepository
	fileRepo  FileRepository
	userRepo  UserRepository
}

func NewRedisRepository(client RedisClient) *RedisRepository {
	return &RedisRepository{
		textRepo:  textrepo.NewRedisRepository(client),
		imageRepo: imagerepo.NewRedisRepository(client),
		fileRepo:  filerepo.NewRedisRepository(client),
		userRepo:  userrepo.NewRedisRepository(client),
	}
}

func (r *RedisRepository) Text() TextRepository {
	return r.textRepo
}

func (r *RedisRepository) Image() ImageRepository {
	return r.imageRepo
}

func (r *RedisRepository) File() FileRepository {
	return r.fileRepo
}

func (r *RedisRepository) User() UserRepository {
	return r.userRepo
}
