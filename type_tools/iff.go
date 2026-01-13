package type_tools

func Ternary[T any](cond bool, t, f T) T {
	if cond {
		return t
	}

	return f
}

func LazyTernary[T any](cond bool, t, f func() T) T {
	if cond {
		return t()
	}

	return f()
}
