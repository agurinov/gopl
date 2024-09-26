package cpy

// #include <Python.h>
import "C"

import (
	"fmt"
	"unsafe"
)

func GetPythonHome() (string, error) {
	wcHome := C.Py_GetPythonHome()
	if wcHome == nil {
		return "", nil
	}

	cHome := C.Py_EncodeLocale(wcHome, nil)
	if cHome == nil {
		return "", fmt.Errorf("fail to call Py_EncodeLocale")
	}

	defer C.PyMem_Free(unsafe.Pointer(cHome))

	return C.GoString(cHome), nil
}
