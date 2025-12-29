package enumerable

type filterImpl[T any] struct {
	Impl
	source Interface[T]
	filter func(T) bool
}

func (f *filterImpl[T]) Next() (next T, ok bool) {
	for {
		next, ok = f.source.Next()
		if !ok {
			return
		}

		if f.filter(next) {
			return
		}
	}
}
