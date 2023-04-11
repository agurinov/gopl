package cpy

// #include <Python.h>
import "C"

import (
	"context"
	"errors"
)

func Initialize() error {
	C.Py_Initialize()

	return nil
}

func Initialized() bool {
	return C.Py_IsInitialized() != 0
}

func SetInterruptFromContext(
	ctx context.Context,
	done <-chan struct{},
) {
	select {
	case <-ctx.Done():
		C.PyErr_SetInterrupt()
	case <-done:
	}
}

func Finalize() error {
	NewGIL().Lock()

	if C.Py_FinalizeEx() == -1 {
		return errors.New("something went wrong with Py_Finalize")
	}

	return nil
}

func GetVersion() string {
	return C.GoString(
		C.Py_GetVersion(),
	)
}

func ThreadsInitialized() bool {
	return C.PyEval_ThreadsInitialized() == 1
}
