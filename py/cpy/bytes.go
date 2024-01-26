package cpy

// #include <Python.h>
import "C"
import "unsafe"

func ToBytesFromString(str string) *Object {
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))

	return C.PyBytes_FromStringAndSize(
		cStr,
		C.Py_ssize_t(len(str)),
	)
}

func ToBytes(b []byte) *Object {
	var (
		bytesPtr unsafe.Pointer
		bytesLen = C.Py_ssize_t(len(b))
	)

	if len(b) != 0 {
		bytesPtr = unsafe.Pointer(&b[0])
	}

	return C.PyBytes_FromStringAndSize(
		(*C.char)(bytesPtr),
		bytesLen,
	)
}
