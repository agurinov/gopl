package cpy

// #include <Python.h>
import "C"

import (
	"fmt"
	"unsafe"
)

func GetPath() (string, error) {
	wcPath := C.Py_GetPath()
	if wcPath == nil {
		return "", nil
	}

	cPath := C.Py_EncodeLocale(wcPath, nil)
	if cPath == nil {
		return "", fmt.Errorf("fail to call Py_EncodeLocale")
	}

	defer C.PyMem_Free(unsafe.Pointer(cPath))

	return C.GoString(cPath), nil
}

func AppendPath(path string) error {
	pyPath, err := GetPath()
	if err != nil {
		return err
	}

	return SetPath(pyPath + ":" + path)
}

func SetPath(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	decodedCPath := C.Py_DecodeLocale(cPath, nil)
	if decodedCPath == nil {
		return fmt.Errorf("fail to call Py_DecodeLocale on '%s'", path)
	}

	C.Py_SetPath(decodedCPath)

	return nil
}
