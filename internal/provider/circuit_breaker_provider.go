package provider

import (
	"github.com/hynek-systems/hynek-poi/internal/circuitbreaker"
	"github.com/hynek-systems/hynek-poi/internal/domain"
)

type CircuitBreakerProvider struct {
	inner Provider
	cb    *circuitbreaker.CircuitBreaker
}

func NewCircuitBreakerProvider(inner Provider, cb *circuitbreaker.CircuitBreaker) *CircuitBreakerProvider {
	return &CircuitBreakerProvider{
		inner: inner,
		cb:    cb,
	}
}

func (p *CircuitBreakerProvider) Name() string {
	return p.inner.Name()
}

func (p *CircuitBreakerProvider) Search(query domain.SearchQuery) ([]domain.POI, error) {

	if !p.cb.Allow() {
		return nil, circuitbreaker.ErrCircuitOpen
	}

	results, err := p.inner.Search(query)

	if err != nil {
		p.cb.Failure()
		return nil, err
	}

	p.cb.Success()

	return results, nil
}
