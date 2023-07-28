package creational_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"

	c "github.com/agurinov/gopl/patterns/creational"
	pl_testing "github.com/agurinov/gopl/testing"
)

var (
	NopWriter    = bytes.NewBufferString("nop")
	UsefulWriter = bytes.NewBufferString("foobar")
)

func TestObject_Construct(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputObject    MyStruct
		inputOptions   []MyStructOption
		expectedObject MyStruct
		pl_testing.TestCase
	}{
		"case00: nil options; without initial": {
			inputObject:    MyStruct{},
			inputOptions:   []MyStructOption{nil, nil, nil},
			expectedObject: MyStruct{},
		},
		"case01: with options; without initial": {
			inputObject: MyStruct{},
			inputOptions: []MyStructOption{
				WithI(100500),
				WithS("foobar"),
				WithW(UsefulWriter),
			},
			expectedObject: MyStruct{i: 100500, s: "foobar", w: UsefulWriter},
		},
		"case02: without options; with initial": {
			inputObject:    MyStruct{i: 999999, s: "foo", w: NopWriter},
			expectedObject: MyStruct{i: 999999, s: "foo", w: NopWriter},
		},
		"case03: with options; with initial; partial overriding": {
			inputObject: MyStruct{i: 999999, s: "foo", w: NopWriter},
			inputOptions: []MyStructOption{
				WithS("foobar"),
				WithW(UsefulWriter),
			},
			expectedObject: MyStruct{i: 999999, s: "foobar", w: UsefulWriter},
		},
		"case04: failed options": {
			inputObject:  MyStruct{},
			inputOptions: []MyStructOption{WithErr(io.EOF)},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			object, err := c.Construct(
				tc.inputObject,
				tc.inputOptions...,
			)

			tc.CheckError(t, err)
			require.Equal(t, tc.expectedObject, object)
		})
	}
}

func TestObject_Extended(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputOptions   []MyStructOptionWithContext
		expectedObject MyStruct
		pl_testing.TestCase
	}{
		"case00: nil options; validation error": {
			inputOptions: []MyStructOptionWithContext{nil, nil, nil},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(validator.ValidationErrors),
			},
		},
		"case01: validation ok; init err": {
			inputOptions: []MyStructOptionWithContext{
				WithICtx(100500),
				WithSCtx("foobar"),
				WithWCtx(UsefulWriter),
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"case02: validation ok; init ok": {
			inputOptions: []MyStructOptionWithContext{
				WithICtx(100),
				WithSCtx("foobar"),
				WithWCtx(UsefulWriter),
			},
			expectedObject: MyStruct{i: 100, s: "foobar", w: UsefulWriter},
		},
		"case03: failed options": {
			inputOptions: []MyStructOptionWithContext{
				WithErrCtx(io.EOF),
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			object, err := NewExtended(context.TODO(), tc.inputOptions...)

			tc.CheckError(t, err)
			require.Equal(t, tc.expectedObject, object)
		})
	}
}
