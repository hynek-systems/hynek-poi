package cache

import (
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

type Cache interface {
	Get(key string) ([]domain.POI, bool)

	Set(key string, value []domain.POI, ttl time.Duration)
}
