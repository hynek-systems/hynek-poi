package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func (c *RedisCache) Client() *redis.Client {
	return c.client
}

func NewRedisCache(addr string, password string, db int) *RedisCache {

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (c *RedisCache) Get(key string) ([]domain.POI, bool) {

	val, err := c.client.Get(c.ctx, key).Result()

	if err != nil {
		return nil, false
	}

	var pois []domain.POI

	err = json.Unmarshal([]byte(val), &pois)

	if err != nil {
		return nil, false
	}

	return pois, true
}

func (c *RedisCache) Set(key string, value []domain.POI, ttl time.Duration) {

	data, err := json.Marshal(value)

	if err != nil {
		return
	}

	c.client.Set(c.ctx, key, data, ttl)
}

var _ Cache = (*RedisCache)(nil)
