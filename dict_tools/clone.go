package dict_tools

func Clone[K comparable, V any](in map[K]V) map[K]V {
	out := make(map[K]V)
	for k, v := range in {
		out[k] = v
	}
	return out
}
