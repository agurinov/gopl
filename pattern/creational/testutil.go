//go:build test_unit

package creational

import (
	"bytes"
	"io"
	"strconv"
)

// Dummy types, stubs and vars to use in testing within nested packages
// This file under test_unit go build tag !!!

const (
	DefaultI = 999999
	DefaultS = "foo"

	CustomI = 100500
	CustomS = "foobar"
)

type (
	MyInterface interface {
		Do() string
	}
	MyStruct struct {
		W io.Writer
		S string
		I int
	}
	MyStructOption = Option[MyStruct]
)

var (
	NopWriter    = bytes.NewBufferString("nop")
	UsefulWriter = bytes.NewBufferString("foobar")
)

var DefaultMyStruct = MyStruct{
	I: DefaultI,
	S: DefaultS,
	W: NopWriter,
}

// Preset of options with success and fail effect
var (
	SuccessObjectOptions = []MyStructOption{
		nil, WithI(CustomI),
		nil, WithS(CustomS),
		nil, WithW(UsefulWriter),
		nil,
	}
	FailObjectOptions = []MyStructOption{
		nil, WithI(CustomI),
		nil, WithS(CustomS),
		nil, WithW(UsefulWriter),
		nil, WithErr(io.EOF),
		nil,
	}
)

// Do implements MyInterface interface
func (o MyStruct) Do() string {
	return strconv.Itoa(o.I) + o.S
}

func WithI(i int) MyStructOption {
	return func(o *MyStruct) error {
		o.I = i

		return nil
	}
}

func WithS(s string) MyStructOption {
	return func(o *MyStruct) error {
		o.S = s

		return nil
	}
}

func WithW(w io.Writer) MyStructOption {
	return func(o *MyStruct) error {
		o.W = w

		return nil
	}
}

func WithErr(err error) MyStructOption {
	return func(_ *MyStruct) error {
		return err
	}
}
