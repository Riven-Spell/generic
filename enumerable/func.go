package enumerable

type funcImpl[O any] struct {
	Impl
	f func() (O, bool)
}

func (f *funcImpl[O]) Next() (O, bool) {
	return f.f()
}
