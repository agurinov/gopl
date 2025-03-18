package x

func SliceOrNil[T any](in []T) any {
	if len(in) == 0 {
		return nil
	}

	return in
}

func ValueOrNil[T comparable](in T) any {
	var zero T

	if in == zero {
		return nil
	}

	return in
}
