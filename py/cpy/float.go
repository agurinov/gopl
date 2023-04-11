package cpy

// #include <Python.h>
import "C"

func ToFloat(f float64) *Object {
	cDouble := C.double(f)

	return C.PyFloat_FromDouble(cDouble)
}

func FromFloat(obj *Object) float64 {
	cDouble := C.PyFloat_AsDouble(obj)

	return float64(cDouble)
}
