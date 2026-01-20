package conc_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"

	"github.com/samix73/game/conc"
)

func TestMap(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	mapFn := func(i int) string {
		return fmt.Sprintf("%v", i)
	}

	input := []int{1, 2, 3, 4, 5}
	output := make([]string, 0, len(input))
	for result := range conc.Map(t.Context(), input, mapFn, conc.WithMaxGoroutines(3)) {
		output = append(output, result)
	}

	for _, item := range input {
		assert.Contains(t, output, mapFn(item))
	}
}
