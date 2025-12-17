package list_tools

import (
	"github.com/Riven-Spell/generic/enumerator"
	"github.com/Riven-Spell/generic/ptr_tools"
)

type NewListEnumeratorOptions struct {
	// Clone defaults to true-- it describes whether or not the list will be cloned to create the enumerator.
	// The list should not change while the enumerator is running if this is false.
	Clone *bool
}

func (o *NewListEnumeratorOptions) defaults() {
	o.Clone = ptr_tools.EnsureNotNil(o.Clone, true)
}

// NewListEnumerator creates a new enumerator[T]
func NewListEnumerator[T any](in []T, opts ...NewListEnumeratorOptions) enumerator.Interface[T] {
	opt := FirstOrZero(opts)
	opt.defaults()

	target := in
	if *opt.Clone {
		target = Clone(in)
	}

	return &listEnum[T]{
		target: target,
	}
}

type listEnum[T any] struct {
	target []T
	index  uint64
}

func (l *listEnum[T]) Reset() {
	l.index = 0
}

func (l *listEnum[T]) Revert() bool {
	if l.index > 0 {
		l.index--
		return true
	}

	return false
}

func (l *listEnum[T]) RemainingReverts() uint64 {
	return l.index
}

func (l *listEnum[T]) HasErrored() error {
	return nil // listEnum can't error
}

func (l *listEnum[T]) HasNext() bool {
	return len(l.target) > 0
}

func (l *listEnum[T]) Next() T {
	out := l.target[0]
	l.target = l.target[1:]
	return out
}
