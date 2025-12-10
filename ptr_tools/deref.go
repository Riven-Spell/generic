package ptr_tools

// DerefOrDefault attempts to dereference in, returning Default if it is nil.
func DerefOrDefault[T any](in *T, Default T) T {
	if in != nil {
		return *in
	}

	return Default
}

// DerefOrZero attempts to dereference in, returning the zero value of T if it is nil.
func DerefOrZero[T any](in *T) (zero T) {
	return DerefOrDefault(in, zero)
}
