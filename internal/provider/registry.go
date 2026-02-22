package provider

import (
	"time"

	"github.com/hynek-systems/hynek-poi/internal/circuitbreaker"
	"github.com/hynek-systems/hynek-poi/internal/config"
)

type RegisteredProvider struct {
	Provider Provider
	Priority int
}

func BuildProviders(cfg config.ProvidersConfig) []RegisteredProvider {

	var result []RegisteredProvider

	// Google
	if cfg.Google.Enabled {

		base := NewGoogleProvider(cfg.Google.ApiKey)

		cb := circuitbreaker.New(
			3,              // failures before open
			30*time.Second, // reset timeout
		)

		protected := NewCircuitBreakerProvider(base, cb)

		result = append(result, RegisteredProvider{
			Provider: protected,
			Priority: cfg.Google.Priority,
		})
	}

	// OSM
	if cfg.OSM.Enabled {

		base := NewOSMProvider()

		cb := circuitbreaker.New(
			3,
			30*time.Second,
		)

		protected := NewCircuitBreakerProvider(base, cb)

		result = append(result, RegisteredProvider{
			Provider: protected,
			Priority: cfg.OSM.Priority,
		})
	}

	return result
}
