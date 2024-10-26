package a

func Filter[E any](s []E, f func(e E) bool) func(func(int, E) bool) {
	return func(yield func(int, E) bool) {
		for i, v := range s {
			if !f(v) {
				continue
			}
			if !yield(i, v) {
				break
			}
		}
	}
}

func UnhandledYieldFilter[E any](s []E, f func(e E) bool) func(func(int, E) bool) {
	return func(yield func(int, E) bool) {
		for i, v := range s {
			if !f(v) {
				continue
			}
			_ = yield(i, v) // want "yield is not handled"
		}
	}
}
