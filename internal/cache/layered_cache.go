package cache

import (
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

type LayeredCache struct {
	l1 Cache
	l2 Cache
}

func NewLayeredCache(l1 Cache, l2 Cache) *LayeredCache {
	return &LayeredCache{
		l1: l1,
		l2: l2,
	}
}

func (c *LayeredCache) Get(key string) ([]domain.POI, bool) {

	// Try L1 (memory)
	if value, found := c.l1.Get(key); found {
		return value, true
	}

	// Try L2 (redis)
	if value, found := c.l2.Get(key); found {

		// populate L1
		c.l1.Set(key, value, 5*time.Minute)

		return value, true
	}

	return nil, false
}

func (c *LayeredCache) Set(key string, value []domain.POI, ttl time.Duration) {

	// Write to both layers
	c.l1.Set(key, value, ttl)
	c.l2.Set(key, value, ttl)
}

var _ Cache = (*LayeredCache)(nil)
