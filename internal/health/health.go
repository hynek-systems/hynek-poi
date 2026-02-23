package health

import (
	"context"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type Checker struct {
	Redis *redis.Client
}

func New(redis *redis.Client) *Checker {
	return &Checker{
		Redis: redis,
	}
}

// Liveness probe
func (c *Checker) HealthHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

// Readiness probe
func (c *Checker) ReadyHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Check Redis
	if err := c.Redis.Ping(ctx).Err(); err != nil {

		http.Error(w, "Redis not ready", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("READY"))
}
