package helpers

func First[S ~[]E, E any](iterator S) (E, bool) {
	for _, item := range iterator {
		return item, true
	}

	var zero E
	return zero, false
}
