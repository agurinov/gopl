package x

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints" //nolint:exptostd
)

const (
	base10  = 10
	bitSize = 64
)

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

	if slices.Contains(empty, in) {
		return zero
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

func SafeAssert[T any](in any) T {
	var zero T

	typed, ok := in.(T)
	if !ok {
		return zero
	}

	return typed
}

func ParseNumber[T constraints.Integer | constraints.Float](in string) (T, error) {
	var zero T

	in = strings.ReplaceAll(in, ",", "")
	in = strings.TrimSpace(in)

	if in == "" {
		return zero, nil
	}

	if i, err := strconv.ParseInt(in, base10, bitSize); err == nil {
		return T(i), nil
	}

	if f, err := strconv.ParseFloat(in, bitSize); err == nil {
		return T(f), nil
	}

	return zero, fmt.Errorf("can't parse number: %s", in)
}
