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

func MapConvert[K1, K2 comparable, V1, V2 any](
	in map[K1]V1,
	mapF func(K1, V1) (K2, V2),
) map[K2]V2 {
	out := make(map[K2]V2, len(in))

	for k, v := range in {
		k2, v2 := mapF(k, v)

		out[k2] = v2
	}

	return out
}

func MapFilter[K comparable, V any](
	in map[K]V,
	useF func(K, V) bool,
) map[K]V {
	out := make(map[K]V, len(in))

	for k, v := range in {
		if useF(k, v) {
			out[k] = v
		}
	}

	return out
}
