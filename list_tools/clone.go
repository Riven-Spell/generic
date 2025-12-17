package list_tools

func Clone[T any](in []T) []T {
	out := make([]T, len(in))
	copy(out, in)
	return out
}
