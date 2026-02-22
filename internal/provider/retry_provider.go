package provider

import (
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

type RetryProvider struct {
	provider Provider
	retries  int
}

func NewRetryProvider(provider Provider, retries int) Provider {

	return &RetryProvider{
		provider: provider,
		retries:  retries,
	}
}

func (p *RetryProvider) Name() string {

	return p.provider.Name()
}

func (p *RetryProvider) Search(query domain.SearchQuery) ([]domain.POI, error) {

	var lastErr error

	for i := 0; i <= p.retries; i++ {

		result, err := p.provider.Search(query)

		if err == nil {
			return result, nil
		}

		lastErr = err

		time.Sleep(100 * time.Millisecond)
	}

	return nil, lastErr
}
