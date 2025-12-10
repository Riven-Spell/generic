package ptr_tools

// To allows converting a explicit value, i.e. a constant to a pointer value.
// This involves a copy, and effectively shorthands the copy-then-reference.
func To[T any](in T) *T {
	return &in
}
