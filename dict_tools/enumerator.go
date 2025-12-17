package dict_tools

type KeyValue[K, V any] struct {
	Key   K
	Value V
}

type mapEnumerator[K comparable, V any] struct {
	idx        uint64
	targetKeys []K
	target     map[K]V
}

func (m mapEnumerator[K, V]) HasNext() bool {
	//TODO implement me
	panic("implement me")
}

func (m mapEnumerator[K, V]) HasErrored() error {
	//TODO implement me
	panic("implement me")
}

func (m mapEnumerator[K, V]) Reset() {
	//TODO implement me
	panic("implement me")
}

func (m mapEnumerator[K, V]) Revert() bool {
	//TODO implement me
	panic("implement me")
}

func (m mapEnumerator[K, V]) RemainingReverts() uint64 {
	//TODO implement me
	panic("implement me")
}

func (m mapEnumerator[K, V]) Next() T {
	//TODO implement me
	panic("implement me")
}
