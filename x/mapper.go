package x

func SliceMap[T1, T2 any](
	in []T1,
	mapF func(T1) T2,
) []T2 {
	out := make([]T2, 0, len(in))

	for i := range in {
		out = append(out, mapF(in[i]))
	}

	return out
}

func SliceMapError[T1, T2 any](
	in []T1,
	mapF func(T1) (T2, error),
) (
	[]T2,
	error,
) {
	if len(in) == 0 {
		return nil, nil
	}

	out := make([]T2, 0, len(in))

	for i := range in {
		m, err := mapF(in[i])
		if err != nil {
			return nil, err
		}

		out = append(out, m)
	}

	return out, nil
}
