package list_tools

func Last[T any](list []T) T {
	return list[len(list)-1]
}
