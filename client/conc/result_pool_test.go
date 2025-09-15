package conc_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"

	"github.com/samix73/game/client/conc"
)

func TestResultPool(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	p := conc.NewResultPool[int](t.Context())

	for i := range 5 {
		p.Go(func(ctx context.Context) (int, error) {
			return i, nil
		})
	}

	results, err := p.Wait()
	assert.NoError(t, err)
	assert.Len(t, results, 5)
}

func TestResultPool_Error(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	p := conc.NewResultPool[int](t.Context())

	for i := range 5 {
		time.Sleep(100 * time.Millisecond)

		p.Go(func(ctx context.Context) (int, error) {
			if i == 2 {
				return 0, assert.AnError
			}

			return i, nil
		})
	}

	results, err := p.Wait()
	assert.Error(t, err)
	assert.Len(t, results, 2)
}
