package x

func MapToSlice[K comparable, V any, E any](
	in map[K]V,
	mapF func(K, V) E,
) []E {
	out := make([]E, 0, len(in))

	for k, v := range in {
		out = append(out, mapF(k, v))
	}

	return out
}

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
