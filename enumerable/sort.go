package enumerable

import (
	"cmp"
	"slices"
)

type sortFuncImpl[T any] struct {
	Impl

	toCollect Interface[T]
	out       Interface[T]
	sortFunc  func(a, b T) int
}

func (s *sortFuncImpl[T]) Next() (T, bool) {
	if s.out == nil {
		list := Collect(s.toCollect)
		slices.SortFunc(list, s.sortFunc)

		s.out = FromList(list, false)
	}

	return s.out.Next()
}

type sortImpl[T cmp.Ordered] struct {
	Impl

	toCollect Interface[T]
	out       Interface[T]
}

func (s *sortImpl[T]) Next() (T, bool) {
	if s.out == nil {
		list := Collect(s.toCollect)
		slices.Sort(list)

		s.out = FromList(list, false)
	}

	return s.out.Next()
}
