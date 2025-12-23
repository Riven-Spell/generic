package enumerable

type listEnumerator[T any] struct {
	Impl

	target []T
}

func (l *listEnumerator[T]) Next() (value T, ok bool) {
	if len(l.target) == 0 {
		return // default 0 false
	}

	ok = true
	value = l.target[0]
	l.target = l.target[1:]
	return
}
