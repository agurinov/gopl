package x

func Unique[T comparable](in []T) []T {
	var (
		set = make(map[T]struct{}, len(in))
		out = make([]T, 0, len(in))
	)

	for i := range in {
		element := in[i]

		if _, exists := set[element]; exists {
			continue
		}

		set[element] = struct{}{}

		out = append(out, element)
	}

	return out
}

func FilterOutEmpty[T comparable](in []T) []T {
	var (
		zero     T
		filtered = make([]T, 0, len(in))
	)

	for i := range in {
		if in[i] == zero {
			continue
		}

		filtered = append(filtered, in[i])
	}

	if len(filtered) == 0 {
		return nil
	}

	return filtered
}
