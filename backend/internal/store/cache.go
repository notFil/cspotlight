package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type RedisStore struct {
	Client *redis.Client
}

func NewRedis(addr string) *RedisStore {
	return &RedisStore{
		Client: redis.NewClient(&redis.Options{
			Addr: addr,
		})}
}

func (r *RedisStore) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisStore) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisStore) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
