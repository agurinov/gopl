package cpy

// #include <Python.h>
import "C"

import (
	"fmt"
	"unsafe"
)

func ImportModule(name string) (*Object, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	module := C.PyImport_ImportModule(cName)
	if module == nil {
		return nil, LastError()
	}

	return module, nil
}

func ImportModuleItem(module *Object, name *Object) (*Object, error) {
	moduleDict := C.PyModule_GetDict(module)
	if moduleDict == nil {
		return nil, LastError()
	}

	item, err := GetDictItem(moduleDict, name)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func GetDictItem(dict *Object, key *Object) (*Object, error) {
	obj := C.PyDict_GetItemWithError(dict, key)
	if obj == nil {
		if err := LastError(); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("python dict key %q doesn't present", key)
	}

	return obj, nil
}

func ImportGlobalDict() (*Object, error) {
	globalDict := C.PyImport_GetModuleDict()
	if globalDict == nil {
		return nil, LastError()
	}

	return globalDict, nil
}
