package x

func MapKeys[K comparable, V any](
	in map[K]V,
) []K {
	keys := make([]K, 0, len(in))

	for k := range in {
		keys = append(keys, k)
	}

	return keys
}

func MapClone[K comparable, V any](
	in map[K]V,
) map[K]V {
	out := make(map[K]V, len(in))

	for k := range in {
		out[k] = in[k]
	}

	return out
}
