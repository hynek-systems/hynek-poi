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

		// timeout
		withTimeout := NewTimeoutProvider(
			base,
			cfg.Google.Timeout,
		)

		// retry
		withRetry := NewRetryProvider(
			withTimeout,
			cfg.Google.Retries,
		)

		// circuit breaker
		cb := circuitbreaker.New(3, 30*time.Second)

		protected := NewCircuitBreakerProvider(
			withRetry,
			cb,
		)

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

		withTimeout := NewTimeoutProvider(
			base,
			cfg.OSM.Timeout,
		)

		// retry
		withRetry := NewRetryProvider(
			withTimeout,
			cfg.OSM.Retries,
		)

		// circuit breaker
		protected := NewCircuitBreakerProvider(
			withRetry,
			cb,
		)

		result = append(result, RegisteredProvider{
			Provider: protected,
			Priority: cfg.OSM.Priority,
		})
	}

	// Foursquare
	if cfg.Foursquare.Enabled {

		base := NewFoursquareProvider(cfg.Foursquare.ApiKey)

		withTimeout := NewTimeoutProvider(
			base,
			cfg.Foursquare.Timeout,
		)

		withRetry := NewRetryProvider(
			withTimeout,
			cfg.Foursquare.Retries,
		)

		cb := circuitbreaker.New(3, 30*time.Second)

		protected := NewCircuitBreakerProvider(
			withRetry,
			cb,
		)

		result = append(result, RegisteredProvider{
			Provider: protected,
			Priority: cfg.Foursquare.Priority,
		})
	}

	return result
}
