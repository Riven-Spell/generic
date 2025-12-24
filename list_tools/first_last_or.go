package list_tools

// FirstOrDefault returns the first element of a list, or the specified default.
func FirstOrDefault[T any](list []T, Default T) T {
	if len(list) > 0 {
		return list[0]
	}

	return Default
}

// FirstOrZero returns the first element of a list, or the specified
func FirstOrZero[T any](list []T) (zero T) {
	return FirstOrDefault(list, zero)
}

func LastOrDefault[T any](list []T, Default T) T {
	if len(list) > 0 {
		return list[len(list)-1]
	}

	return Default
}

func LastOrZero[T any](list []T) (zero T) {
	return LastOrDefault(list, zero)
}
