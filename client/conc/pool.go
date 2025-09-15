package conc

import (
	"context"
	"sync"

	"github.com/samix73/game/conc/panics"
)

type opt struct {
	maxGoroutines int
	panicHandler  func(recovered *panics.Recovered)
}

type Options func(*opt)

// WithMaxGoroutines sets the maximum number of goroutines that can be
// running concurrently in the pool.
//
// If the value is negative, it will panic.
// If the value is zero, it will not limit the number of goroutines.
func WithMaxGoroutines(n int) Options {
	return func(o *opt) {
		if n < 0 {
			panic("max goroutines in a pool cannot be negative")
		}

		o.maxGoroutines = n
	}
}

func WithPanicHandler(handler func(recovered *panics.Recovered)) Options {
	return func(o *opt) {
		if handler == nil {
			panic("panic handler cannot be nil")
		}

		o.panicHandler = handler
	}
}

type limiter chan struct{}

func (l limiter) release() {
	if l != nil {
		<-l
	}
}

func (l limiter) acquire() {
	if l != nil {
		l <- struct{}{}
	}
}

// Pool is a goroutine pool that allows you to run functions concurrently
// with a limit on the number of goroutines that can run at the same time.
type Pool struct {
	ctx          context.Context
	cancel       context.CancelFunc
	panicHandler func(recovered *panics.Recovered)

	wg      sync.WaitGroup
	limiter limiter
}

// NewPool creates a new Pool with the provided options.
func NewPool(ctx context.Context, o ...Options) *Pool {
	var options opt
	for _, opt := range o {
		opt(&options)
	}

	ctx, cancel := context.WithCancel(ctx)
	p := &Pool{
		ctx:          ctx,
		cancel:       cancel,
		panicHandler: options.panicHandler,
	}

	if options.maxGoroutines > 0 {
		p.limiter = make(limiter, options.maxGoroutines)
	}

	return p
}

// Go runs the provided function in a goroutine. If the maximum number of
// goroutines is reached, it will block until a goroutine becomes available.
// In case of a panic, it will recover from the panic, call the panic handler
// if set, and cancel the pool context, which will stop all other goroutines from running.
// If the pool context is already cancelled, it will not run the function.
func (p *Pool) Go(fn func(ctx context.Context)) {
	p.wg.Add(1)

	p.limiter.acquire()

	go p.worker(fn)
}

func (p *Pool) worker(fn func(ctx context.Context)) {
	defer func() {
		p.limiter.release()
		p.wg.Done()
	}()

	select {
	case <-p.ctx.Done():
		return
	default:
	}

	recovered := panics.Try(func() {
		fn(p.ctx)
	})
	if recovered != nil {
		if p.panicHandler == nil {
			recovered.RePanic()
		}

		p.panicHandler(recovered)

		p.cancel()
	}
}

// MaxGoroutines returns the maximum number of goroutines that can run concurrently in the pool.
func (p *Pool) MaxGoroutines() int {
	if p.limiter == nil {
		return 0
	}

	return cap(p.limiter)
}

// Wait waits for all goroutines in the pool to finish.
func (p *Pool) Wait() {
	p.wg.Wait()
	p.cancel()
}
