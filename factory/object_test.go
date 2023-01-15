//go:build test_unit

package pl_factory_test

import (
	"io"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	pl_factory "github.com/agurinov/gopl.git/factory"
	pl_testing "github.com/agurinov/gopl.git/testing"
)

type (
	MyInterface interface {
		Do() string
	}
	MyStruct struct {
		i int
		s string
	}
	MyStructOption = pl_factory.Option[MyStruct]
)

var (
	New       = pl_factory.NewObject[MyStruct]
	MustNew   = pl_factory.MustNewObject[MyStruct]
	NewAs     = func(opts ...MyStructOption) (MyInterface, error) { return New(opts...) }
	MustNewAs = func(opts ...MyStructOption) MyInterface { return MustNew(opts...) }
)

// Catch Must methods with panics as errors
var (
	recoveredMustNew = func(opts ...MyStructOption) (object MyStruct, err error) {
		defer func() {
			if r := recover(); r != nil {
				err, _ = r.(error)
			}
		}()

		object = MustNew(opts...)

		return
	}
	recoveredMustNewAs = func(opts ...MyStructOption) (object MyInterface, err error) {
		defer func() {
			if r := recover(); r != nil {
				err, _ = r.(error)
			}
		}()

		object = MustNewAs(opts...)

		return
	}
)

// Preset of options with success and fail effect
var (
	successOptions = []MyStructOption{
		nil, WithI(100500),
		nil, WithS("foobar"),
		nil,
	}
	failOptions = []MyStructOption{
		nil, WithI(100500),
		nil, WithS("foobar"),
		nil, WithErr(io.EOF),
		nil,
	}
)

func (o MyStruct) Do() string {
	return strconv.Itoa(o.i) + o.s
}

func WithI(i int) MyStructOption {
	return func(o *MyStruct) error {
		o.i = i

		return nil
	}
}

func WithS(s string) MyStructOption {
	return func(o *MyStruct) error {
		o.s = s

		return nil
	}
}

func WithErr(err error) MyStructOption {
	return func(_ *MyStruct) error {
		return err
	}
}

func TestConstructor(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputOptions   []MyStructOption
		expectedStruct MyStruct
		expectedIface  MyInterface
		pl_testing.TestCase
	}{
		"fail": {
			inputOptions: failOptions,
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"success with options": {
			inputOptions:   successOptions,
			expectedStruct: MyStruct{i: 100500, s: "foobar"},
			expectedIface:  MyInterface(MyStruct{i: 100500, s: "foobar"}),
		},
		"success without options": {
			expectedStruct: MyStruct{},
			expectedIface:  MyInterface(MyStruct{}),
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			Struct1, struct1Err := New(tc.inputOptions...)
			Struct2, struct2Err := recoveredMustNew(tc.inputOptions...)
			Iface1, iface1Err := NewAs(tc.inputOptions...)
			Iface2, iface2Err := recoveredMustNewAs(tc.inputOptions...)

			tc.CheckError(t, struct1Err)
			tc.CheckError(t, struct2Err)
			tc.CheckError(t, iface1Err)
			tc.CheckError(t, iface2Err)

			require.Equal(t, tc.expectedStruct, Struct1)
			require.Equal(t, tc.expectedStruct, Struct2)
			require.Equal(t, tc.expectedIface, Iface1)
			require.Equal(t, tc.expectedIface, Iface2)
		})
	}
}
