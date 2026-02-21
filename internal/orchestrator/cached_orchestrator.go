package orchestrator

import (
	"time"

	"github.com/hynek-systems/hynek-poi/internal/cache"
	"github.com/hynek-systems/hynek-poi/internal/domain"
)

type CachedOrchestrator struct {
	inner Orchestrator
	cache cache.Cache
	ttl   time.Duration
}

func NewCached(inner Orchestrator, cache cache.Cache, ttl time.Duration) *CachedOrchestrator {
	return &CachedOrchestrator{
		inner: inner,
		cache: cache,
		ttl:   ttl,
	}
}

func (c *CachedOrchestrator) Search(query domain.SearchQuery) ([]domain.POI, error) {

	key := cache.BuildKey(query)

	// cache hit
	if cached, found := c.cache.Get(key); found {
		return cached, nil
	}

	// cache miss
	results, err := c.inner.Search(query)

	if err != nil {
		return nil, err
	}

	c.cache.Set(key, results, c.ttl)

	return results, nil
}
