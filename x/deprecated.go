package x

import "cmp"

// Deprecated: Use SliceConvert instead.
//
//go:fix inline
func SliceMap[T1, T2 any](
	in []T1,
	mapF func(T1) T2,
) []T2 {
	return SliceConvert(in, mapF)
}

// Deprecated: Use SliceConvertError instead.
//
//go:fix inline
func SliceMapError[T1, T2 any](
	in []T1,
	mapF func(T1) (T2, error),
) (
	[]T2,
	error,
) {
	return SliceConvertError(in, mapF)
}

// Deprecated: Use cmp.Or instead.
//
//go:fix inline
func Coalesce[T comparable](in ...T) T {
	return cmp.Or(in...)
}
