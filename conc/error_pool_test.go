package conc_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"

	"gitlab.com/abios/gohelpers/conc"
)

func TestErrorPool(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	p := conc.NewErrorPool(t.Context())

	var executed int64
	for i := range 5 {
		time.Sleep(100 * time.Millisecond)

		p.Go(func(ctx context.Context) error {
			atomic.AddInt64(&executed, 1)

			if i == 2 {
				return assert.AnError
			}

			return nil
		})
	}

	err := p.Wait()
	assert.Error(t, err)
	assert.Equal(t, int64(3), executed)
}
