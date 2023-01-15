//go:build test_unit

package creational_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	c "github.com/agurinov/gopl.git/pattern/creational"
	pl_testing "github.com/agurinov/gopl.git/testing"
)

var (
	// Struct constructors
	New     = c.NewObject[c.MyStruct]
	MustNew = c.MustNewObject[c.MyStruct]

	// Polymorphic constructors
	NewAs     = func(opts ...c.MyStructOption) (c.MyInterface, error) { return New(opts...) }
	MustNewAs = func(opts ...c.MyStructOption) c.MyInterface { return MustNew(opts...) }
)

// Catch Must methods with panics as errors
var (
	recoveredMustNew = func(opts ...c.MyStructOption) (object c.MyStruct, err error) {
		defer func() {
			if r := recover(); r != nil {
				err, _ = r.(error)
			}
		}()

		object = MustNew(opts...)

		return
	}
	recoveredMustNewAs = func(opts ...c.MyStructOption) (object c.MyInterface, err error) {
		defer func() {
			if r := recover(); r != nil {
				err, _ = r.(error)
			}
		}()

		object = MustNewAs(opts...)

		return
	}
)

func TestObject_New_MustNew(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputOptions   []c.MyStructOption
		expectedObject c.MyStruct
		expectedIface  c.MyInterface
		pl_testing.TestCase
	}{
		"fail": {
			inputOptions: c.FailObjectOptions,
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"success with options": {
			inputOptions:   c.SuccessObjectOptions,
			expectedObject: c.MyStruct{I: 100500, S: "foobar", W: c.UsefulWriter},
			expectedIface:  c.MyInterface(c.MyStruct{I: 100500, S: "foobar", W: c.UsefulWriter}),
		},
		"success without options": {
			expectedObject: c.MyStruct{},
			expectedIface:  c.MyInterface(c.MyStruct{}),
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

			require.Equal(t, tc.expectedObject, Struct1)
			require.Equal(t, tc.expectedObject, Struct2)
			require.Equal(t, tc.expectedIface, Iface1)
			require.Equal(t, tc.expectedIface, Iface2)
		})
	}
}

func TestObject_Construct(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputObject    c.MyStruct
		inputOptions   []c.MyStructOption
		expectedObject c.MyStruct
		pl_testing.TestCase
	}{
		"without options and without initial": {
			inputObject:    c.MyStruct{},
			expectedObject: c.MyStruct{},
		},
		"with options and without initial": {
			inputObject:    c.MyStruct{},
			inputOptions:   c.SuccessObjectOptions,
			expectedObject: c.MyStruct{I: 100500, S: "foobar", W: c.UsefulWriter},
		},
		"without options and with initial": {
			inputObject:    c.DefaultMyStruct,
			expectedObject: c.DefaultMyStruct,
		},
		"with options and with initial (partial overriding)": {
			inputObject: c.DefaultMyStruct,
			inputOptions: []c.MyStructOption{
				c.WithS("foobar"),
				c.WithW(c.UsefulWriter),
			},
			expectedObject: c.MyStruct{I: 999999, S: "foobar", W: c.UsefulWriter},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			Struct, err := c.ConstructObject(
				tc.inputObject,
				tc.inputOptions...,
			)

			tc.CheckError(t, err)

			require.Equal(t, tc.expectedObject, Struct)
		})
	}
}
