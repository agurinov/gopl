package cpy

// #include <Python.h>
import "C"

func Tuple(args ...*Object) (*Object, error) {
	tuple := C.PyTuple_New(
		C.Py_ssize_t(len(args)),
	)
	if tuple == nil {
		return nil, LastError()
	}

	// Each PyTuple_SetItem steals reference to underlying Object
	for i := range args {
		if C.PyTuple_SetItem(tuple, C.Py_ssize_t(i), args[i]) != 0 {
			return nil, LastError()
		}
	}

	return tuple, nil
}

func FromTuple(tuple *Object) ([]*Object, error) {
	args := make([]*Object, C.PyTuple_Size(tuple))

	for i := range args {
		item := C.PyTuple_GetItem(tuple, C.Py_ssize_t(i))
		if item == nil {
			return nil, LastError()
		}

		args[i] = item
	}

	return args, nil
}
