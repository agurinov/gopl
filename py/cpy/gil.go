package cpy

// #include <Python.h>
import "C"

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type gil struct {
	state atomic.Uintptr
}

func (g *gil) Lock() {
	runtime.LockOSThread()

	g.state.Store(
		uintptr(C.PyGILState_Ensure()),
	)
}

func (g *gil) Unlock() {
	C.PyGILState_Release(
		C.PyGILState_STATE(g.state.Load()),
	)

	runtime.UnlockOSThread()
}

func NewGIL() sync.Locker {
	return new(gil)
}

func SaveThread() *C.PyThreadState {
	return C.PyEval_SaveThread()
}
