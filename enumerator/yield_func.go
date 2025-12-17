package enumerator

type NewYielderOptions struct {
	// BufferSize is the revert buffer available in the yield enumerator.
	// This will determine the amount of terms that can be rolled back.
	BufferSize *uint
}

func (o *NewYielderOptions) defaults() {

}

func NewYielder[T any](f func(<-chan T)) Interface[T] {
	outChannel := make(chan T)
	errChannel := make(chan error)
	ctrlChannel := make(chan int)

	/*
		The control channel is designed to assist in
	*/
}
