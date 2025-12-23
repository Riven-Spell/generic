package enumerable

type flattenImpl[T any] struct {
	Impl

	source  Interface[[]T]
	current []T
}

func (f *flattenImpl[T]) Next() (out T, ok bool) {
	for f.source == nil || len(f.current) == 0 {
		f.current, ok = f.source.Next()
		if !ok {
			return
		}
	}

	out = f.current[0]
	f.current = f.current[1:]

	if len(f.current) == 0 {
		f.current = nil
	}

	return out, true
}
