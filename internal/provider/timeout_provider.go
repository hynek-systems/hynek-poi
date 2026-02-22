package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
)

type TimeoutProvider struct {
	provider Provider
	timeout  time.Duration
}

func NewTimeoutProvider(provider Provider, timeout time.Duration) Provider {

	return &TimeoutProvider{
		provider: provider,
		timeout:  timeout,
	}
}

func (p *TimeoutProvider) Name() string {

	return p.provider.Name()
}

func (p *TimeoutProvider) Search(query domain.SearchQuery) ([]domain.POI, error) {

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	resultChan := make(chan []domain.POI, 1)
	errorChan := make(chan error, 1)

	go func() {

		result, err := p.provider.Search(query)

		if err != nil {
			errorChan <- err
			return
		}

		resultChan <- result
	}()

	select {

	case result := <-resultChan:
		return result, nil

	case err := <-errorChan:
		return nil, err

	case <-ctx.Done():
		return nil, fmt.Errorf("provider timeout: %s", p.provider.Name())
	}
}
