package cpy

// #include <Python.h>
import "C"

type Object = C.PyObject

func (obj *Object) IncRef() { C.Py_IncRef(obj) }
func (obj *Object) DecRef() { C.Py_DecRef(obj) }

func (obj *Object) String() string {
	return FromString(obj)
}

func (obj *Object) SetAttr(k, v *Object) error {
	res := C.PyObject_SetAttr(obj, k, v)

	if res != 0 {
		return LastError()
	}

	return nil
}

func (obj *Object) GetAttr(k *Object) (*Object, error) {
	v := C.PyObject_GetAttr(obj, k)

	if v == nil {
		return nil, LastError()
	}

	return v, nil
}
