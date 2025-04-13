package x

func Keys[K comparable, V any](
	in map[K]V,
) []K {
	keys := make([]K, 0, len(in))

	for k := range in {
		keys = append(keys, k)
	}

	return keys
}
