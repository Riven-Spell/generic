package ptr_tools

func EnsureWithZero[T any](in *T) *T {
	var Zero T
	return EnsureWithDefault(in, Zero)
}

func EnsureWithDefault[T any](in *T, Default T) *T {
	if in == nil {
		return &Default
	}

	return in
}
