//go:build test_unit

package pl_factory_test

import (
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	pl_factory "github.com/agurinov/gopl.git/factory"
	pl_testing "github.com/agurinov/gopl.git/testing"
)

type (
	Foo struct {
		a int
		b int
	}
	FooOption  = pl_factory.Option[Foo]
	FooFactory = pl_factory.Factory[Foo]
)

var (
	New        = pl_factory.NewObject[Foo]
	MustNew    = pl_factory.MustNewObject[Foo]
	NewFactory = pl_factory.New[Foo]
)

var (
	WithA = func(a int) FooOption {
		return func(f *Foo) error {
			f.a = a

			return nil
		}
	}
	WithB = func(b int) FooOption {
		return func(f *Foo) error {
			f.b = b

			return nil
		}
	}
	WithErr = func(err error) FooOption {
		return func(_ *Foo) error {
			return err
		}
	}
)

func TestFactory_NewObject(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		inputOptions   []FooOption
		expectedObject Foo
		pl_testing.TestCase
	}{
		"fail": {
			inputOptions: []FooOption{
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
			inputOptions: []FooOption{
				WithA(100),
				WithB(500),
			},
			expectedObject: Foo{a: 100, b: 500},
		},
		"success without options": {},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			factory := NewFactory(tc.inputOptions...)

			// Catch MustNew method with panics.
			factoryMustNewObject := func() (Foo, error) {
				var (
					wg     sync.WaitGroup
					object Foo
					err    error
				)

				wg.Add(1)
				go func() {
					defer wg.Done()
					defer func() {
						if r := recover(); r != nil {
							err, _ = r.(error)
						}
					}()

					object = factory.MustNewObject()
				}()
				wg.Wait()

				return object, err
			}

			object1, object1Err := factory.NewObject()
			tc.CheckError(t, object1Err)

			object2, object2Err := factoryMustNewObject()
			tc.CheckError(t, object2Err)

			object3, object3Err := New(tc.inputOptions...)
			tc.CheckError(t, object3Err)

			if !tc.MustFail {
				require.Equal(t, tc.expectedObject, object1)
				require.Equal(t, tc.expectedObject, object2)
				require.Equal(t, tc.expectedObject, object3)
			}
		})
	}
}
