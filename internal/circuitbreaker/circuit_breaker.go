package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker struct {
	mu sync.Mutex

	state State

	failures    int
	lastFailure time.Time

	failureThreshold int
	openTimeout      time.Duration
}

func New(failureThreshold int, openTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: failureThreshold,
		openTimeout:      openTimeout,
	}
}

func (cb *CircuitBreaker) Allow() bool {

	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {

	case StateClosed:
		return true

	case StateOpen:

		if time.Since(cb.lastFailure) > cb.openTimeout {
			cb.state = StateHalfOpen
			return true
		}

		return false

	case StateHalfOpen:
		return true

	default:
		return false
	}
}

func (cb *CircuitBreaker) Success() {

	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures = 0
	cb.state = StateClosed
}

func (cb *CircuitBreaker) Failure() {

	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.failures >= cb.failureThreshold {
		cb.state = StateOpen
	}
}

var ErrCircuitOpen = errors.New("circuit breaker is open")
