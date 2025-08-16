package conc

import (
	"context"
	"sync/atomic"
)

// ErrorPool is a goroutine pool that allows you to run functions concurrently
// with a limit on the number of goroutines that can run at the same time.
type ErrorPool struct {
	pool  *Pool
	error atomic.Value
}

// NewErrorPool creates a new Pool with the provided options.
func NewErrorPool(ctx context.Context, o ...Options) *ErrorPool {
	return &ErrorPool{
		pool: NewPool(ctx, o...),
	}
}

// Go runs the provided function in a goroutine. If the maximum number of
// goroutines is reached, it will block until a goroutine becomes available.
// Go captures the first error returned by the function and cancels the pool context,
// which will stop all other goroutines from running. If the function panics,
// it will recover from the panic, store the recovered value as an error, and cancel the pool context.
// If the pool context is already cancelled, it will not run the function.
func (p *ErrorPool) Go(fn func(ctx context.Context) error) {
	p.pool.Go(func(ctx context.Context) {
		if err := fn(ctx); err != nil {
			p.error.CompareAndSwap(nil, err)
			p.pool.cancel()
		}
	})
}

// MaxGoroutines returns the maximum number of goroutines that can run concurrently in the pool.
func (p *ErrorPool) MaxGoroutines() int {
	return p.pool.MaxGoroutines()
}

// Wait waits for all goroutines in the pool to finish. It will return the
// first error encountered by any of the goroutines, or nil if no errors occurred.
func (p *ErrorPool) Wait() error {
	p.pool.Wait()

	if err := p.error.Load(); err != nil {
		return err.(error)
	}

	return nil
}
