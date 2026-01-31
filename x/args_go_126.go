//go:build go1.26

package x

// Deprecated: Use new(T) instead.
//
//go:fix inline
func Ptr[T any](in T) *T {
	return new(in)
}
