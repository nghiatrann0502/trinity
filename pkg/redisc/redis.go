package redisc

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisEngine interface {
	Health() (bool, error)
	GetRedis() *redis.Client
	Close()
}

type redisCache struct {
	redis *redis.Client
}

var _ RedisEngine = (*redisCache)(nil)

func NewRedisClient(url string) (RedisEngine, error) {
	opts, err := redis.ParseURL(string(url))
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	return &redisCache{
		redis: client,
	}, nil
}

func (r *redisCache) Health() (bool, error) {
	ctx := context.Background()
	_, err := r.redis.Ping(ctx).Result()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *redisCache) GetRedis() *redis.Client {
	return r.redis
}

func (r *redisCache) Close() {
	if r.redis != nil {
		r.redis.Close()
	}
}
