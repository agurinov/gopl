package cpy

// #include <Python.h>
import "C"

func ToInt(i int) *Object {
	cLong := C.long(i)

	return C.PyLong_FromLong(cLong)
}

func FromInt(obj *Object) int {
	cLong := C.PyLong_AsLong(obj)

	return int(cLong)
}
