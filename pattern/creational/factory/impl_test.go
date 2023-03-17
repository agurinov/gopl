//go:build test_unit

package factory_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	c "github.com/agurinov/gopl/pattern/creational"
	"github.com/agurinov/gopl/pattern/creational/factory"
	pl_testing "github.com/agurinov/gopl/testing"
)

type (
	MyStructFactory = c.Factory[c.MyStruct]
)

var (
	successFactory MyStructFactory = factory.MustNew[c.MyStruct](
		factory.WithInitialObject(c.DefaultMyStruct),
		factory.WithOptions(c.SuccessObjectOptions...),
	)
	failFactory MyStructFactory = factory.MustNew[c.MyStruct](
		factory.WithInitialObject(c.DefaultMyStruct),
		factory.WithOptions(c.FailObjectOptions...),
	)
)

func TestFactory_NewObject(t *testing.T) {
	pl_testing.Init(t)

	onlyInitialFactory, err := factory.New[c.MyStruct](
		factory.WithInitialObject(c.DefaultMyStruct),
	)
	require.NoError(t, err)

	emptyObjectFactory, err := factory.New[c.MyStruct]()
	require.NoError(t, err)

	cases := map[string]struct {
		inputFactory   MyStructFactory
		expectedObject c.MyStruct
		pl_testing.TestCase
	}{
		"fail": {
			inputFactory: failFactory,
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: io.EOF,
			},
		},
		"success with options with initial": {
			inputFactory:   successFactory,
			expectedObject: c.MyStruct{I: 100500, S: "foobar", W: c.UsefulWriter},
		},
		"success without options with initial": {
			inputFactory:   onlyInitialFactory,
			expectedObject: c.DefaultMyStruct,
		},
		"success without options without initial": {
			inputFactory:   emptyObjectFactory,
			expectedObject: c.MyStruct{},
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

			require.Equal(t, tc.expectedObject, Struct1)
			require.Equal(t, tc.expectedObject, Struct2)
			require.NotSame(t, Struct1, Struct2)
		})
	}
}
