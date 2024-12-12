package resilience

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/eapache/go-resiliency/breaker"
	"github.com/eapache/go-resiliency/deadline"
	"github.com/eapache/go-resiliency/retrier"
)

var (
	ErrTimeout     = deadline.ErrTimedOut
	ErrCircuitOpen = breaker.ErrBreakerOpen
	ErrMaxRetries  = errors.New("maximum retries exceeded")
)

type ResilienceService interface {
	ExecuteWithAll(ctx context.Context, name string, op func() error) error
}

type resilienceService struct {
	breakers  map[string]*breaker.Breaker
	retriers  map[string]*retrier.Retrier
	deadlines map[string]*deadline.Deadline
	mu        sync.RWMutex
	config    Config
}

type Config struct {
	// Circuit Breaker settings
	BreakerFailures  int           // Number of failures before opening
	BreakerResetTime time.Duration // Time before attempting reset

	// Retrier settings
	RetryAttempts int           // Maximum number of retry attempts
	RetryBackoff  time.Duration // Initial backoff duration

	// Deadline settings
	DefaultTimeout time.Duration // Default timeout for operations
}

func DefaultConfig() Config {
	return Config{
		BreakerFailures:  5,
		BreakerResetTime: 30 * time.Second,
		RetryAttempts:    3,
		RetryBackoff:     100 * time.Millisecond,
		DefaultTimeout:   5 * time.Second,
	}
}

func NewResilienceService(config Config) *resilienceService {
	return &resilienceService{
		breakers:  make(map[string]*breaker.Breaker),
		retriers:  make(map[string]*retrier.Retrier),
		deadlines: make(map[string]*deadline.Deadline),
		config:    config,
	}
}

// ExecuteWithRetries performs an operation with retry logic
func (s *resilienceService) ExecuteWithRetries(name string, op func() error) error {
	r := s.getOrCreateRetrier(name)
	return r.Run(op)
}

// ExecuteWithBreaker performs an operation with circuit breaker protection
func (s *resilienceService) ExecuteWithBreaker(name string, op func() error) error {
	b := s.getOrCreateBreaker(name)
	return b.Run(op)
}

// ExecuteWithDeadline performs an operation with a deadline
func (s *resilienceService) ExecuteWithDeadline(name string, timeout time.Duration, op func(stopper <-chan struct{}) error) error {
	d := s.getOrCreateDeadline(name, timeout)
	return d.Run(op)
}

// ExecuteWithAll combines all resilience patterns
func (s *resilienceService) ExecuteWithAll(ctx context.Context, name string, op func() error) error {
	b := s.getOrCreateBreaker(name)
	r := s.getOrCreateRetrier(name)
	d := s.getOrCreateDeadline(name, s.config.DefaultTimeout)

	done := make(chan error, 1)

	// Run the operation with all protections
	err := d.Run(func(stopper <-chan struct{}) error {
		go func() {
			err := b.Run(func() error {
				return r.Run(op)
			})
			done <- err
		}()

		select {
		case err := <-done:
			return err
		case <-stopper:
			return ErrTimeout
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	return err
}

// Helper methods to get or create components
func (s *resilienceService) getOrCreateBreaker(name string) *breaker.Breaker {
	s.mu.Lock()
	defer s.mu.Unlock()

	if b, exists := s.breakers[name]; exists {
		return b
	}

	b := breaker.New(s.config.BreakerFailures, 1, s.config.BreakerResetTime)
	s.breakers[name] = b
	return b
}

func (s *resilienceService) getOrCreateRetrier(name string) *retrier.Retrier {
	s.mu.Lock()
	defer s.mu.Unlock()

	if r, exists := s.retriers[name]; exists {
		return r
	}

	r := retrier.New(retrier.ExponentialBackoff(s.config.RetryAttempts, s.config.RetryBackoff), nil)
	s.retriers[name] = r
	return r
}

func (s *resilienceService) getOrCreateDeadline(name string, timeout time.Duration) *deadline.Deadline {
	s.mu.Lock()
	defer s.mu.Unlock()

	if d, exists := s.deadlines[name]; exists {
		return d
	}

	d := deadline.New(timeout)
	s.deadlines[name] = d
	return d
}

// Reset resets all resilience components for a given name
func (s *resilienceService) Reset(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.breakers, name)
	delete(s.retriers, name)
	delete(s.deadlines, name)
}
