package cache

import (
	"sync"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

type cacheItem struct {
	value      []domain.POI
	expiration time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]cacheItem),
	}
}

func (c *MemoryCache) Get(key string) ([]domain.POI, bool) {

	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

func (c *MemoryCache) Set(key string, value []domain.POI, ttl time.Duration) {

	c.mu.Lock()

	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}

	c.mu.Unlock()
}
