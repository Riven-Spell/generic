package dict_tools

func Keys[K comparable, V any](in map[K]V) []K {
	out := make([]K, len(in))
	idx := 0

	for k, _ := range in {
		out[idx] = k
	}

	return out
}

func Values[K comparable, V any](in map[K]V) []V {
	out := make([]V, len(in))
	idx := 0

	for _, v := range in {
		out[idx] = v
	}

	return out
}
