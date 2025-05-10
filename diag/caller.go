package diag

import (
	"fmt"
	"runtime"
	"strings"
)

//nolint:gomnd,mnd
func CallerName() string {
	pc, _, _, _ := runtime.Caller(1) //nolint:dogsled

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
			packageName = parts[n-3]
			structName  = strings.Trim(parts[n-2], "()")
			methodName  = parts[n-1]
		)

		return fmt.Sprintf("%s.%s.%s", packageName, structName, methodName)
	case n == 2:
		var (
			packageName  = parts[n-2]
			functionName = parts[n-1]
		)

		return fmt.Sprintf("%s.%s", packageName, functionName)
	default:
		return "unsplittable"
	}
}
