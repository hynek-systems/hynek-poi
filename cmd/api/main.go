package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/cache"
	"github.com/hynek-systems/hynek-poi/internal/circuitbreaker"
	"github.com/hynek-systems/hynek-poi/internal/config"
	"github.com/hynek-systems/hynek-poi/internal/domain"
	"github.com/hynek-systems/hynek-poi/internal/metrics"
	"github.com/hynek-systems/hynek-poi/internal/orchestrator"
	"github.com/hynek-systems/hynek-poi/internal/provider"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var orch *orchestrator.CachedOrchestrator

func searchHandler(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	metrics.RequestsTotal.WithLabelValues("/v1/search").Inc()

	defer func() {
		metrics.RequestDuration.
			WithLabelValues("/v1/search").
			Observe(time.Since(start).Seconds())
	}()

	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lng, _ := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)

	query := domain.SearchQuery{
		Latitude:  lat,
		Longitude: lng,
		Radius:    1000,
		Limit:     50,
	}

	results, err := orch.Search(query)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(results)
}

func main() {

	cfg := config.Load()

	metrics.Register()

	osmProvider := provider.NewOSMProvider()

	cb := circuitbreaker.New(3, 30*time.Second)

	protectedProvider := provider.NewCircuitBreakerProvider(
		osmProvider,
		cb,
	)

	weighted := orchestrator.NewWeighted([]orchestrator.ProviderConfig{
		{
			Provider: protectedProvider,
			Weight:   10,
		},
	})

	memoryCache := cache.NewMemoryCache()

	redisCache := cache.NewRedisCache(
		cfg.Redis.Addr,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)

	layeredCache := cache.NewLayeredCache(
		memoryCache,
		redisCache,
	)

	orch = orchestrator.NewCached(
		weighted,
		layeredCache,
		cfg.Cache.TTL,
	)

	http.HandleFunc("/v1/search", searchHandler)

	http.Handle("/metrics", promhttp.Handler())

	addr := ":" + strconv.Itoa(cfg.Server.Port)

	log.Println("Hynek POI listening on", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
