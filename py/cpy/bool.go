package cpy

// #include <Python.h>
import "C"

var (
	False = C.Py_False
	True  = C.Py_True
)

func ToBool(b bool) *Object {
	if b {
		return True
	}

	return False
}

func FromBool(obj *Object) bool {
	return obj == True
}
