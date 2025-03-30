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

func EmptyIf[T comparable](in T, empty ...T) T {
	var zero T

	for i := range empty {
		if in == empty[i] {
			return zero
		}
	}

	return in
}

func Ptr[T any](in T) *T {
	return &in
}

func FromPtr[T any](in *T) T {
	var zero T

	if in == nil {
		return zero
	}

	return *in
}
