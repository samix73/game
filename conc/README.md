# Conc
A utility package that will provide common golang concurrency patterns.

## Pool

Pool is a pool of goroutines used to execute tasks concurrently.

### Example

Running a set of tasks and waiting for them to finish without limiting the number of goroutines:
```go
package main

import (
	"context"

	"gitlab.com/abios/gohelpers/conc"
)

func main() {
	ctx := context.Background()

	pool := conc.NewPool(ctx)
	
	for range 10 {
		pool.Go(func(ctx context.Context) {
			// Do some work here
		})	
    }
	
	pool.Wait()
}
```

Running a set of tasks and waiting for them to finish with a limited number of goroutines:
```go
package main

import (
	"context"

	"gitlab.com/abios/gohelpers/conc"
)

func main() {
	ctx := context.Background()

	pool := conc.NewPool(ctx, conc.WithMaxGoroutines(3))
	
	for range 10 {
		// Pool will not spawn more than 3 goroutines at a time
        // If more tasks are submitted, Go will block until a goroutine is available
		pool.Go(func(ctx context.Context) {
			// Do some work here
		})	
    }
	
	pool.Wait()
}
```

Run some tasks and get the aggregated results:
```go
package main

import (
	"context"

	"gitlab.com/abios/gohelpers/conc"
)

func main() {
	ctx := context.Background()

	pool := conc.NewResultPool[int](ctx)
	
	for i := range 10 {
		pool.Go(func(ctx context.Context) (int, error) {
			// Do some work here
			
			return i, nil
		})
    }
	
	results, err := pool.Wait()
	if err != nil {
        // Handle error
    }
	
	// ... process results
}
```

Run tasks and cancel them if any of them fails:
```go
package main

import (
	"context"
	"errors"

	"gitlab.com/abios/gohelpers/conc"
)

func main() {
	ctx := context.Background()

	pool := conc.NewErrorPool(ctx)
	
	for range 10 {
		pool.Go(func(ctx context.Context) error {
			// Do some work here
			
			return errors.New("some error") 
		})	
    }
	
	if err := pool.Wait(); err != nil {
		// Handle error
    }
}
```

Panic recovery is also supported.

```go
package main

import (
	"context"

	"gitlab.com/abios/gohelpers/conc"
	"gitlab.com/abios/gohelpers/conc/panics"
)

func main() {
	ctx := context.Background()

	panicHandler := func(r *panics.Recovered) {
		// Handle panic here
	}

	pool := conc.NewPool(ctx, conc.WithPanicHandler(panicHandler))

	for range 10 {
		pool.Go(func(ctx context.Context) {
			// Do some work here

			panic("some panic")
		})
	}

	pool.Wait()
}
```


## Iterators

Iterators are a way to iterate over a set of items concurrently.

Map:
```go
package main


import (
    "context"
	
    "gitlab.com/abios/gohelpers/conc"
)

func main() {
	ctx := context.Background()
	
	mapFn := func(ctx context.Context, item int) (int, error) {
        return item * 2, nil
	}
	
	for result := range conc.Map(ctx, []int{1, 2, 3}, mapFn, conc.WithMaxGoroutines(2)) {
		// Process result
    }
}

```