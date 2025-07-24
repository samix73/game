package helpers

import "iter"

func First[T any](iterator iter.Seq[T]) (T, bool) {
	for item := range iterator {
		return item, true
	}

	var zero T
	return zero, false
}
