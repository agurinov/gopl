package x

import (
	"errors"
)

func FlattenChans[T any](chs ...<-chan T) []T {
	out := make([]T, 0, len(chs))

	for _, ch := range chs {
		for range cap(ch) {
			select {
			case obj := <-ch:
				out = append(out, obj)
			default:
			}
		}
	}

	return out
}

func FlattenErrors(chs ...<-chan error) error {
	errs := FlattenChans(chs...)

	return errors.Join(errs...)
}
