//go:build test_unit

package py_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/agurinov/gopl/py"
	"github.com/agurinov/gopl/py/cpy"
	pl_testing "github.com/agurinov/gopl/testing"
)

func pyWork(
	ctx context.Context,
	pyModuleName string,
	pyFunctionName string,
	mapperToPython py.MapperToPython,
	mapperFromPython py.MapperFromPython,
) error {
	interpreter, err := py.NewInterpreter(
		py.WithPythonPath(
			"testdata",
			filepath.Join("testdata/pypkg", py.VenvPath("venv")),
		),
		py.WithEnsureGIL(true),
		py.WithEnsureVersion("3.8.16"),
	)
	if err != nil {
		return err
	}

	fn, err := interpreter.GetCallable(
		ctx,
		pyModuleName+"."+pyFunctionName,
	)
	if err != nil {
		return err
	}

	if err := fn.Call(ctx, mapperToPython, mapperFromPython); err != nil {
		return err
	}

	return nil
}

func TestInterpreter_PyFunc_InverseFloat(t *testing.T) {
	pl_testing.Init(t)

	const (
		pyModuleName   = "pypkg"
		pyFunctionName = "inverse_float"
	)

	cases := map[string]struct {
		inputFloat       float64
		expectedResponse float64
		pl_testing.TestCase
	}{
		"0": {inputFloat: -0.145, expectedResponse: 0.145},
		"1": {inputFloat: 1.2, expectedResponse: -1.2},
		"2": {inputFloat: 100500, expectedResponse: -100500},
		"3": {inputFloat: -100.500, expectedResponse: 100.500},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			ctx := context.TODO()

			var response float64

			mapperToPython := func() ([]*py.Object, error) {
				return []*py.Object{
					cpy.ToFloat(tc.inputFloat),
				}, nil
			}
			mapperFromPython := func(res *py.Object) error {
				response = cpy.FromFloat(res)

				return nil
			}

			tc.CheckError(t, pyWork(
				ctx,
				pyModuleName,
				pyFunctionName,
				mapperToPython,
				mapperFromPython,
			))
			require.Equal(t, tc.expectedResponse, response)
		})
	}
}

func TestInterpreter_PyFunc_RaiseException(t *testing.T) {
	pl_testing.Init(t)

	const (
		pyModuleName   = "pypkg"
		pyFunctionName = "raise_exception"
	)

	cases := map[string]struct {
		inputString string
		pl_testing.TestCase
	}{
		"0": {
			inputString: "foo",
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(cpy.ExceptionError),
			},
		},
		"1": {
			inputString: "foo",
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(cpy.ExceptionError),
			},
		},
		"2": {
			inputString: "foo",
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(cpy.ExceptionError),
			},
		},
		"3": {
			inputString: "foo",
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(cpy.ExceptionError),
			},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			ctx := context.TODO()

			mapperToPython := func() ([]*py.Object, error) {
				return []*py.Object{
					cpy.ToString(tc.inputString),
				}, nil
			}

			tc.CheckError(t, pyWork(
				ctx,
				pyModuleName,
				pyFunctionName,
				mapperToPython,
				nil,
			))
		})
	}
}

func TestInterpreter_PyFunc_HelloWorld(t *testing.T) {
	pl_testing.Init(t)

	const (
		pyModuleName   = "pypkg"
		pyFunctionName = "hello_world"
	)

	cases := map[string]struct {
		pl_testing.TestCase
	}{
		"0": {},
		"1": {},
		"2": {},
		"3": {},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			ctx := context.TODO()

			tc.CheckError(t, pyWork(
				ctx,
				pyModuleName,
				pyFunctionName,
				nil, nil,
			))
		})
	}
}

func TestInterpreter_PyFunc_InverseComplex(t *testing.T) {
	pl_testing.Init(t)

	const (
		pyModuleName   = "pypkg"
		pyFunctionName = "inverse_complex"
	)

	cases := map[string]struct {
		pl_testing.TestCase
		inputString    string
		expectedString string
		inputFloat     float64
		inputInt       int
		expectedFloat  float64
		expectedInt    int
		inputBool      bool
		expectedBool   bool
	}{
		"0": {
			inputString: "foobar", inputBool: true, inputFloat: 1.2345, inputInt: -999,
			expectedString: "raboof", expectedBool: false, expectedFloat: -1.2345, expectedInt: 999,
		},
		"1": {
			inputString: "lolkek", inputBool: false, inputFloat: -100.500, inputInt: 35,
			expectedString: "keklol", expectedBool: true, expectedFloat: 100.500, expectedInt: -35,
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			ctx := context.TODO()

			var (
				expectedString string
				expectedBool   bool
				expectedFloat  float64
				expectedInt    int
			)

			mapperToPython := func() ([]*py.Object, error) {
				return []*py.Object{
					cpy.ToString(tc.inputString),
					cpy.ToBool(tc.inputBool),
					cpy.ToFloat(tc.inputFloat),
					cpy.ToInt(tc.inputInt),
				}, nil
			}
			mapperFromPython := func(res *py.Object) error {
				args, err := cpy.FromTuple(res)
				require.NoError(t, err)
				require.Equal(t, 4, len(args))
				expectedString = cpy.FromString(args[0])
				expectedBool = cpy.FromBool(args[1])
				expectedFloat = cpy.FromFloat(args[2])
				expectedInt = cpy.FromInt(args[3])

				return nil
			}

			tc.CheckError(t, pyWork(
				ctx,
				pyModuleName,
				pyFunctionName,
				mapperToPython,
				mapperFromPython,
			))
			require.Equal(t, tc.expectedString, expectedString)
			require.Equal(t, tc.expectedBool, expectedBool)
			require.Equal(t, tc.expectedFloat, expectedFloat)
			require.Equal(t, tc.expectedInt, expectedInt)
		})
	}
}

func TestInterpreter_PyFunc_Contract(t *testing.T) {
	pl_testing.Init(t)

	const (
		pyModuleName   = "pypkg"
		pyFunctionName = "contract"
	)

	type (
		PythonContractRequest struct {
			S string
			B bool
			F float64
			I int
		}
		PythonContractResponse struct {
			MlResult float64
		}
	)

	cases := map[string]struct {
		inputPythonRequest     PythonContractRequest
		expectedPythonResponse PythonContractResponse
		pl_testing.TestCase
	}{
		"s=a": {
			inputPythonRequest:     PythonContractRequest{S: "a"},
			expectedPythonResponse: PythonContractResponse{MlResult: 0.001},
		},
		"b=True": {
			inputPythonRequest:     PythonContractRequest{B: true},
			expectedPythonResponse: PythonContractResponse{MlResult: 0.555},
		},
		"f=100.500": {
			inputPythonRequest:     PythonContractRequest{F: 100.500},
			expectedPythonResponse: PythonContractResponse{MlResult: 100.500},
		},
		"i=100": {
			inputPythonRequest:     PythonContractRequest{I: 100},
			expectedPythonResponse: PythonContractResponse{MlResult: 1.123},
		},
		"undefined": {
			inputPythonRequest:     PythonContractRequest{},
			expectedPythonResponse: PythonContractResponse{MlResult: -999},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			ctx := context.TODO()

			var response PythonContractResponse

			mapperToPython := func() ([]*py.Object, error) {
				module, err := cpy.ImportModule(pyModuleName)
				if err != nil {
					return nil, err
				}

				classObject, err := cpy.ImportModuleItem(module, cpy.ToString("ContractRequest"))
				if err != nil {
					return nil, err
				}

				initArgsTuple, err := cpy.Tuple(
					cpy.ToString(tc.inputPythonRequest.S),
					cpy.ToBool(tc.inputPythonRequest.B),
					cpy.ToFloat(tc.inputPythonRequest.F),
					cpy.ToInt(tc.inputPythonRequest.I),
				)
				if err != nil {
					return nil, err
				}

				classInstance, err := cpy.CallObject(classObject, initArgsTuple)
				if err != nil {
					return nil, err
				}

				return []*py.Object{classInstance}, nil
			}

			mapperFromPython := func(res *py.Object) error {
				v, err := res.GetAttr(cpy.ToString("ml_result"))
				if err != nil {
					return err
				}

				response.MlResult = cpy.FromFloat(v)

				return nil
			}

			tc.CheckError(t, pyWork(
				ctx,
				pyModuleName,
				pyFunctionName,
				mapperToPython,
				mapperFromPython,
			))
			require.Equal(t, tc.expectedPythonResponse, response)
		})
	}
}

func TestInterpreter_PyFunc_PIL(t *testing.T) {
	pl_testing.Init(t)

	const (
		pyModuleName   = "pypkg"
		pyFunctionName = "pil"
	)

	imageBytes, err := os.ReadFile("testdata/meme_1024_576.jpg")
	require.NoError(t, err)
	require.NotNil(t, imageBytes)
	require.Equal(t, 44080, len(imageBytes))

	cases := map[string]struct {
		inputBytes    []byte
		expectedFloat float64
		pl_testing.TestCase
	}{
		"nil blob": {
			inputBytes: nil,
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(cpy.ExceptionError),
			},
		},
		"some static blob - not an image": {
			inputBytes: []byte("foobar"),
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailAsErr: new(cpy.ExceptionError),
			},
		},
		"real image": {
			inputBytes:    imageBytes,
			expectedFloat: 44080,
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			ctx := context.TODO()

			var response float64

			mapperToPython := func() ([]*py.Object, error) {
				return []*py.Object{
					cpy.ToBytes(tc.inputBytes),
				}, nil
			}
			mapperFromPython := func(res *py.Object) error {
				response = cpy.FromFloat(res)

				return nil
			}

			tc.CheckError(t, pyWork(
				ctx,
				pyModuleName,
				pyFunctionName,
				mapperToPython,
				mapperFromPython,
			))
			require.Equal(t, tc.expectedFloat, response)
		})
	}
}

func TestInterpreter_PyFunc_InverseBools(t *testing.T) {
	pl_testing.Init(t)

	const (
		pyModuleName   = "pypkg"
		pyFunctionName = "inverse_bools"
	)

	cases := map[string]struct {
		inputString      string
		inputBool        bool
		expectedResponse bool
		pl_testing.TestCase
	}{
		"0": {inputString: "as-is", inputBool: true, expectedResponse: true},
		"1": {inputString: "as-is", inputBool: false, expectedResponse: false},
		"2": {inputString: "reverse", inputBool: true, expectedResponse: false},
		"3": {inputString: "reverse", inputBool: false, expectedResponse: true},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			ctx := context.TODO()

			var response bool

			mapperToPython := func() ([]*py.Object, error) {
				return []*py.Object{
					cpy.ToString(tc.inputString),
					cpy.ToBool(tc.inputBool),
				}, nil
			}
			mapperFromPython := func(res *py.Object) error {
				response = cpy.FromBool(res)

				return nil
			}

			tc.CheckError(t, pyWork(
				ctx,
				pyModuleName,
				pyFunctionName,
				mapperToPython,
				mapperFromPython,
			))
			require.Equal(t, tc.expectedResponse, response)
		})
	}
}

func TestInterpreter_PyFunc_TimeSleep(t *testing.T) {
	pl_testing.Init(t)

	const (
		pyModuleName   = "time"
		pyFunctionName = "sleep"
	)

	cases := map[string]struct {
		pyModuleName   string
		pyFunctionName string
		inputSeconds   float64
		pl_testing.TestCase
	}{
		"0":  {inputSeconds: 1.0},
		"1":  {inputSeconds: 1.1},
		"2":  {inputSeconds: 1.2},
		"3":  {inputSeconds: 1.3},
		"4":  {inputSeconds: 1.4},
		"5":  {inputSeconds: 1.5},
		"6":  {inputSeconds: 1.6},
		"7":  {inputSeconds: 1.7},
		"8":  {inputSeconds: 1.8},
		"9":  {inputSeconds: 1.9},
		"10": {inputSeconds: 2.0},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			ctx := context.TODO()

			mapperToPython := func() ([]*py.Object, error) {
				return []*py.Object{
					cpy.ToFloat(tc.inputSeconds),
				}, nil
			}

			tc.CheckError(t, pyWork(
				ctx,
				pyModuleName,
				pyFunctionName,
				mapperToPython,
				nil,
			))
		})
	}
}
