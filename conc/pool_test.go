package conc_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"

	"gitlab.com/abios/gohelpers/conc"
	"gitlab.com/abios/gohelpers/conc/panics"
)

func TestPool(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	p := conc.NewPool(t.Context())

	p.Go(func(ctx context.Context) {
		time.Sleep(100 * time.Millisecond)
	})

	p.Go(func(ctx context.Context) {
		time.Sleep(200 * time.Millisecond)
	})

	p.Wait()
}

func TestPool_WithMaxGoroutines(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	const maxGoroutines = 3

	p := conc.NewPool(t.Context(), conc.WithMaxGoroutines(maxGoroutines))

	var running int64
	var expectedMaxConcurrent int64

	for range 30 {
		p.Go(func(ctx context.Context) {
			current := atomic.AddInt64(&running, 1)
			defer atomic.AddInt64(&running, -1)

			if current > expectedMaxConcurrent {
				atomic.StoreInt64(&expectedMaxConcurrent, current)
			}

			time.Sleep(100 * time.Millisecond)
		})
	}

	p.Wait()

	assert.Equal(t, int64(0), running)
	assert.GreaterOrEqual(t, int64(maxGoroutines), expectedMaxConcurrent)
}

func TestPool_WithPanicHandler(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	handlerCalled := atomic.Bool{}
	panicHandler := conc.WithPanicHandler(func(recovered *panics.Recovered) {
		assert.NotNil(t, recovered)
		assert.Equal(t, assert.AnError, recovered.Value)

		assert.True(t, handlerCalled.CompareAndSwap(false, true))
	})

	p := conc.NewPool(t.Context(), panicHandler)

	const toExecute = int64(5)
	var executed int64

	for i := range toExecute {
		time.Sleep(100 * time.Millisecond)

		p.Go(func(ctx context.Context) {
			atomic.AddInt64(&executed, 1)

			if i == 1 {
				panic(assert.AnError)
			}
		})
	}

	p.Wait()

	assert.Equal(t, int64(2), executed)
	assert.True(t, handlerCalled.Load())
}
