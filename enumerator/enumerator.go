package enumerator

type Enumerator[T any] interface {
	HasNext() bool
	Next() T
}
