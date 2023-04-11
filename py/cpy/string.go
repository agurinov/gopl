package cpy

// #include <Python.h>
import "C"
import "unsafe"

func ToString(str string) *Object {
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))

	return C.PyUnicode_FromString(cStr)
}

func FromStringWithError(pyObj *Object) (string, error) {
	pyStr := C.PyObject_Str(pyObj)
	if pyStr == nil {
		return "", LastError()
	}

	defer pyStr.DecRef()

	return C.GoString(
		C.PyUnicode_AsUTF8(pyStr),
	), nil
}

func FromString(pyObj *Object) string {
	s, err := FromStringWithError(pyObj)
	if err != nil {
		return ""
	}

	return s
}
