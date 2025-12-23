package enumerable

type mapImpl[I, O any] struct {
	Impl

	i Interface[I]
	f func(I) O
}

func (m mapImpl[I, O]) Next() (out O, ok bool) {
	var interim I
	interim, ok = m.i.Next()
	if !ok {
		return
	}

	return m.f(interim), ok
}
