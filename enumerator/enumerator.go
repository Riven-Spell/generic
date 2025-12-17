package enumerator

type Interface[T any] interface {
	HasNext() bool
	HasErrored() error

	Reset()                   // reset to step 0
	Revert() bool             // revert to the last step (returns if was possible
	RemainingReverts() uint64 // How far back can we revert?

	Next() T
}
