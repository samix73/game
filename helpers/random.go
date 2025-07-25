package helpers

import (
	"math"
	"math/rand/v2"
)

var (
	randomIntIndex int
	randPool       []int
)

func init() {
	randPool = make([]int, 250_000)
	for i := range randPool {
		randPool[i] = rand.IntN(math.MaxInt)
	}
}

func RandomInt(min, max int) int {
	if randomIntIndex >= len(randPool) {
		randomIntIndex = 0
	}

	value := randPool[randomIntIndex%len(randPool)]
	randomIntIndex++

	return value%(max-min+1) + min
}
