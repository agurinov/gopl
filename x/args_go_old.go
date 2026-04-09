//go:build !go1.26

package x

func Ptr[T any](in T) *T {
	return &in
}
