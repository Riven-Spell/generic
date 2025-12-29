package enumerable

import (
	"github.com/Riven-Spell/generic/dict_tools"
	"github.com/Riven-Spell/generic/list_tools"
)

// interface.go should serve as the interface to the enumerable package, a hub to discover the rest of the package via.
// Unfortunately, because the semantics of methods and generics,
// it is impossible to have one function perfectly chain onto another. i.e. There can never be a FromFunc().Zip() or whatnot,
// which means it mostly has to be done backwards... Zip(Map(FromList(x), func()...), y). This is a huge mess, and it's recommended to just do it line by line.

// RawInterface serves for interacting with the Interface at a type-free level.
// It is not used within the package currently, but is provided as a convenience to end users.
type RawInterface interface {
	impl()
}

// Interface defines the very What the package is interacting with, an enumeration of items.
type Interface[T any] interface {
	impl()
	// Next returns true so long as *this* returned value is valid. False indicates the Interface is contractually dead.
	Next() (T, bool)
}

// FromFunc creates a new Interface off the back of a operator function, which will generate new values.
func FromFunc[O any](f func() (O, bool)) Interface[O] {
	return &funcImpl[O]{
		Impl{},
		f,
	}
}

// FromList creates a new Interface that emits all items in order from the list.
func FromList[O any](in []O, clone bool) Interface[O] {
	if clone {
		in = list_tools.Clone(in)
	}

	return FromFunc(func() (out O, ok bool) {
		for _, v := range in {
			out = v
			ok = true
			break
		}

		return
	})
}

// FromMap creates a new Interface that emits all key value pairs in a semi-random order
// (depending on whichever order Golang returns them)
func FromMap[A comparable, B any](src map[A]B) Interface[AB[A, B]] {
	src = dict_tools.Clone(src)

	return FromFunc(func() (out AB[A, B], ok bool) {
		for k, v := range src {
			delete(src, k)
			out.A, out.B = k, v
			ok = true
			break
		}

		return
	})
}

// Zip makes no assurances about the length of either Interface.
// Zip will return false when either A or B returns false.
// Zip aligns two lists from their 0th index.
func Zip[A, B any](Left Interface[A], Right Interface[B]) Interface[AB[A, B]] {
	return &zipImpl[A, B]{
		Impl{},
		Left, Right,
	}
}

// Map converts an input value to an output value.
func Map[I, O any](Source Interface[I], operator func(I) O) Interface[O] {
	return &mapImpl[I, O]{
		Impl{},
		Source, operator,
	}
}

// Filter returns only the values in which operator returns true to.
func Filter[I any](Source Interface[I], operator func(I) bool) Interface[I] {
	return &filterImpl[I]{
		Impl{},
		Source, operator,
	}
}

// Flatten reduces a incoming set of []Ts down to a single set of Ts.
func Flatten[T any](Source Interface[[]T]) Interface[T] {
	return &flattenImpl[T]{
		Impl{},
		Source, nil,
	}
}

// ForEach fully executes an incoming Interface returning no values, and performing the requested operation on each value
func ForEach[T any](Source Interface[T], do func(T)) {
	Sum[T, any](Source, func(t T, _ any) any {
		do(t)
		return nil
	})
}

// Sum reduces an incoming Interface down to a single value, passing the next item and
func Sum[I, O any](Source Interface[I], operator func(I, O) O) O {
	var out O
	for {
		in, ok := Source.Next()
		if !ok {
			return out
		}

		out = operator(in, out)
	}
}

// Collect reduces an incoming interface
func Collect[O any](source Interface[O]) []O {
	return Sum[O, []O](source, func(object O, list []O) []O {
		return append(list, object)
	})
}

func CollectMap[A comparable, B any](Source Interface[AB[A, B]]) map[A]B {
	return Sum[AB[A, B], map[A]B](Source, func(a AB[A, B], m map[A]B) map[A]B {
		return dict_tools.SafeAdd(m, a.A, a.B)
	})
}
