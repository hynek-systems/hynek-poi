package provider

import (
	"time"

	"github.com/hynek-systems/hynek-poi/internal/circuitbreaker"
	"github.com/hynek-systems/hynek-poi/internal/config"
)

func BuildProviders(cfg config.ProvidersConfig) []Provider {

	var providers []Provider

	if cfg.OSM.Enabled {

		osm := NewOSMProvider()

		cb := circuitbreaker.New(3, 30*time.Second)

		protected := NewCircuitBreakerProvider(osm, cb)

		providers = append(providers, protected)
	}

	// future providers:
	// if cfg.Google.Enabled { ... }

	return providers
}
