package conc

import (
	"context"
)

// Map applies the provided function to each element of the slice
// and returns a sequence of results. The function is executed in parallel
// using a goroutine pool.
func Map[Slice ~[]E, E any, R any](ctx context.Context, slice Slice, fn func(E) R, ops ...Options) <-chan R {
	if len(slice) == 0 {
		return nil
	}

	p := NewPool(ctx, ops...)

	results := make(chan R, len(slice))
	for _, item := range slice {
		p.Go(func(ctx context.Context) {
			results <- fn(item)
		})
	}

	go func() {
		p.Wait()
		close(results)
	}()

	return results
}
