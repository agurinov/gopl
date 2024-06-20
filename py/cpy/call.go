package cpy

// #include <Python.h>
import "C"

func CallObject(fn *Object, args *Object) (*Object, error) {
	res := C.PyObject_CallObject(fn, args)
	if res == nil {
		return nil, LastError()
	}

	return res, nil
}

func IsCallable(obj *Object) bool {
	return C.PyCallable_Check(obj) == 1
}
