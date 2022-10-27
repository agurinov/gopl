//go:build test_unit

package pl_option_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	pl_option "github.com/agurinov/gopl.git/option"
	pl_testing "github.com/agurinov/gopl.git/testing"
)

type (
	SomeStruct struct {
		a int
		b int
	}
	SomeStructOption = pl_option.Option[SomeStruct]
)

var (
	WithA = func(a int) SomeStructOption {
		return func(t *SomeStruct) error {
			t.a = a

			return nil
		}
	}
	WithB = func(b int) SomeStructOption {
		return func(t *SomeStruct) error {
			t.b = b

			return nil
		}
	}
	WithErr = func(err error) SomeStructOption {
		return func(t *SomeStruct) error {
			return err
		}
	}
)

func TestOption_SomeStruct(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputOptions       []SomeStructOption
		expectedSomeStruct SomeStruct
		pl_testing.TestCase
	}{
		"fail": {
			inputOptions: []SomeStructOption{
				WithA(1),
				nil,
				WithB(5),
				WithErr(io.EOF),
			},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"success with options": {
			inputOptions: []SomeStructOption{
				WithA(100),
				WithB(500),
			},
			expectedSomeStruct: SomeStruct{a: 100, b: 500},
		},
		"success without options": {},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			someStruct, err := pl_option.New[SomeStruct](tc.inputOptions...)
			tc.CheckError(t, err)

			if !tc.MustFail {
				require.Equal(t, tc.expectedSomeStruct, someStruct)
			}
		})
	}
}
