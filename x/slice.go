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

func Last[E any](s []E) E {
	var zero E

	if len(s) == 0 {
		return zero
	}

	return s[len(s)-1]
}

func First[E any](s []E) E {
	var zero E

	if len(s) == 0 {
		return zero
	}

	return s[0]
}

func SliceFilter[T any](
	in []T,
	useF func(T) bool,
) []T {
	out := make([]T, 0, len(in))

	for i := range in {
		if useF(in[i]) {
			out = append(out, in[i])
		}
	}

	return out
}

func SliceToMap[K comparable, V any, E any](
	in []E,
	mapF func(E) (K, V),
) map[K]V {
	out := make(map[K]V, len(in))

	for i := range in {
		k, v := mapF(in[i])
		out[k] = v
	}

	return out
}

func SliceConvert[T1, T2 any](
	in []T1,
	mapF func(T1) T2,
) []T2 {
	out := make([]T2, 0, len(in))

	for i := range in {
		out = append(out, mapF(in[i]))
	}

	return out
}

func SliceConvertError[T1, T2 any](
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

func Paginate[T any](
	in []T,
	limit uint,
	offset uint,
) []T {
	if limit == 0 {
		return nil
	}

	sliceLen := uint(len(in))

	if offset >= sliceLen {
		return nil
	}

	end := min(
		offset+limit,
		sliceLen,
	)

	return in[offset:end]
}

func Flatten[T any](
	in [][]T,
) []T {
	if in == nil {
		return nil
	}

	flattened := make([]T, 0, len(in))

	for _, subSlice := range in {
		flattened = append(flattened, subSlice...)
	}

	return flattened
}

func SliceBatch[T any](
	in []T,
	batchSize uint,
) [][]T {
	switch {
	case len(in) == 0:
		return nil
	case batchSize == 0:
		return nil
	}

	var (
		sliceLen   = uint(len(in))
		batchCount = (sliceLen + batchSize - 1) / batchSize
		batches    = make([][]T, 0, batchCount)
	)

	for start := uint(0); start < sliceLen; start += batchSize {
		end := min(
			start+batchSize,
			sliceLen,
		)

		batch := make([]T, end-start)
		copy(batch, in[start:end])

		batches = append(batches, batch)
	}

	return batches
}
