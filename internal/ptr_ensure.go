package internal

// EnsureNotNil returns in, or Default if in is nil.
func EnsureNotNil[T any](in *T, Default T) *T {
	if in == nil {
		return &Default
	}

	return in
}
