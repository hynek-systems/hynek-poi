package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/cache"
	"github.com/hynek-systems/hynek-poi/internal/config"
	"github.com/hynek-systems/hynek-poi/internal/domain"
	"github.com/hynek-systems/hynek-poi/internal/health"
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

	categoriesParam := r.URL.Query().Get("categories")

	var categories []string

	if categoriesParam != "" {
		categories = strings.Split(categoriesParam, ",")
	}

	query := domain.SearchQuery{
		Latitude:   lat,
		Longitude:  lng,
		Radius:     1000,
		Limit:      50,
		Categories: categories,
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

	providers := provider.BuildProviders(cfg.Providers)

	parallel := orchestrator.NewParallel(
		providers,
		3*time.Second,
	)

	memoryCache := cache.NewMemoryCache()

	redisCache := cache.NewRedisCache(
		cfg.Redis.Addr,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)

	healthChecker := health.New(redisCache.Client())

	layeredCache := cache.NewLayeredCache(
		memoryCache,
		redisCache,
	)

	orch = orchestrator.NewCached(
		parallel,
		layeredCache,
		cfg.Cache.TTL,
	)

	http.HandleFunc("/v1/search", searchHandler)
	http.HandleFunc("/health", healthChecker.HealthHandler)
	http.HandleFunc("/ready", healthChecker.ReadyHandler)

	http.Handle("/metrics", promhttp.Handler())

	addr := ":" + strconv.Itoa(cfg.Server.Port)

	log.Println("Hynek POI listening on", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}
