//go:build test_unit

package pl_factory_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	pl_factory "github.com/agurinov/gopl.git/factory"
	pl_testing "github.com/agurinov/gopl.git/testing"
)

var (
	successSingleton MyStructFactory = pl_factory.NewSingleton[MyStruct](successOptions...)
	failSingleton    MyStructFactory = pl_factory.NewSingleton[MyStruct](failOptions...)
)

func TestFactory_Singleton_NewObject(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputFactory   MyStructFactory
		expectedStruct MyStruct
		pl_testing.TestCase
	}{
		"fail": {
			inputFactory: failSingleton,
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"success": {
			inputFactory:   successSingleton,
			expectedStruct: MyStruct{i: 100500, s: "foobar"},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			Struct1, struct1Err := tc.inputFactory.NewObject()
			Struct2, struct2Err := tc.inputFactory.NewObject()

			tc.CheckError(t, struct1Err)
			tc.CheckError(t, struct2Err)

			require.Equal(t, tc.expectedStruct, Struct1)
			require.Equal(t, tc.expectedStruct, Struct2)
			// require.Same(t, Struct1, Struct2)
		})
	}
}