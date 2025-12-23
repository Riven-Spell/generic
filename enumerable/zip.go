package enumerable

type zipImpl[A any, B any] struct {
	Impl

	Left  Interface[A]
	Right Interface[B]
}

func (z *zipImpl[A, B]) Next() (AB[A, B], bool) {
	a, ok := z.Left.Next()
	b, bOk := z.Right.Next()

	return AB[A, B]{a, b}, ok && bOk
}
