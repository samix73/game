package conc

import (
	"context"
	"sync"
)

// ResultPool is a pool that runs functions in goroutines and collects their results.
type ResultPool[T any] struct {
	pool    *ErrorPool
	results []T
	mu      sync.Mutex
}

// NewResultPool creates a new ResultPool with the provided context and options.
func NewResultPool[T any](ctx context.Context, o ...Options) *ResultPool[T] {
	p := &ResultPool[T]{
		pool:    NewErrorPool(ctx, o...),
		results: make([]T, 0),
		mu:      sync.Mutex{},
	}

	return p
}

// Go runs the provided function in a goroutine. If the maximum number of
// goroutines is reached, it will block until a goroutine becomes available.
// If the pool context is already cancelled, it will not run the function.
// If function returns an error, the value will not be added to the results slice.
func (p *ResultPool[T]) Go(fn func(ctx context.Context) (T, error)) {
	p.pool.Go(func(ctx context.Context) error {
		result, err := fn(ctx)
		if err != nil {
			return err
		}

		p.mu.Lock()
		p.results = append(p.results, result)
		p.mu.Unlock()

		return nil
	})
}

func (p *ResultPool[T]) MaxGoroutines() int {
	return p.pool.MaxGoroutines()
}

// Wait waits for all goroutines to finish and returns the collected results.
func (p *ResultPool[T]) Wait() ([]T, error) {
	if err := p.pool.Wait(); err != nil {
		return p.results, err
	}

	return p.results, nil
}
