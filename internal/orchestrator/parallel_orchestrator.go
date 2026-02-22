package orchestrator

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hynek-systems/hynek-poi/internal/domain"
	"github.com/hynek-systems/hynek-poi/internal/provider"
)

type ParallelOrchestrator struct {
	providers []provider.Provider
	timeout   time.Duration
}

func NewParallel(
	providers []provider.Provider,
	timeout time.Duration,
) *ParallelOrchestrator {

	return &ParallelOrchestrator{
		providers: providers,
		timeout:   timeout,
	}
}

func (o *ParallelOrchestrator) Search(query domain.SearchQuery) ([]domain.POI, error) {

	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	resultChan := make(chan []domain.POI, 1)
	errorChan := make(chan error, len(o.providers))

	var wg sync.WaitGroup

	for _, p := range o.providers {

		wg.Add(1)

		go func(provider provider.Provider) {

			defer wg.Done()

			results, err := provider.Search(query)

			if err != nil {
				errorChan <- err
				return
			}

			if len(results) > 0 {

				select {

				case resultChan <- results:
					cancel()

				case <-ctx.Done():
				}
			}

		}(p)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	select {

	case results := <-resultChan:
		return results, nil

	case <-ctx.Done():
		return nil, errors.New("all providers failed or timeout")
	}
}

var _ Orchestrator = (*ParallelOrchestrator)(nil)
