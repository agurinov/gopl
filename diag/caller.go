package diag

import (
	"cmp"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func FunctionName(f any) string {
	v := reflect.ValueOf(f)

	isNil := cmp.Or(
		f == nil,
	)

	isInvalid := cmp.Or(
		!v.IsValid(),
		v.Kind() != reflect.Func,
	)

	switch {
	case isNil:
		return "nil"
	case isInvalid:
		return "invalid"
	}

	pc := v.Pointer()

	return pcName(pc)
}

//nolint:gomnd
func CallerName(skip int) string {
	skip = cmp.Or(skip, 1)

	pc, _, _, _ := runtime.Caller(skip) //nolint:dogsled

	return pcName(pc)
}

//nolint:mnd
func pcName(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	fullName := fn.Name()

	lastSlash := strings.LastIndex(fullName, "/")
	if lastSlash >= 0 {
		fullName = fullName[lastSlash+1:]
	}

	parts := strings.Split(fullName, ".")

	switch n := len(parts); {
	case n < 2:
		return fullName
	case n > 2:
		var (
			packageName = parts[0]
			structName  = strings.Trim(parts[n-2], "()")
			methodName  = strings.TrimSuffix(parts[n-1], "-fm")
		)

		return fmt.Sprintf("%s.%s.%s", packageName, structName, methodName)
	case n == 2:
		var (
			packageName  = parts[0]
			functionName = parts[n-1]
		)

		return fmt.Sprintf("%s.%s", packageName, functionName)
	default:
		return "unsplittable"
	}
}
