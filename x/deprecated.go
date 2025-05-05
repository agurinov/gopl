package x

// Deprecated: Use SliceConvert instead.
func SliceMap[T1, T2 any](
	in []T1,
	mapF func(T1) T2,
) []T2 {
	return SliceConvert(in, mapF)
}

// Deprecated: Use SliceConvertError instead.
func SliceMapError[T1, T2 any](
	in []T1,
	mapF func(T1) (T2, error),
) (
	[]T2,
	error,
) {
	return SliceConvertError(in, mapF)
}
