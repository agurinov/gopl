package cpy

// #include <Python.h>
import "C"
import "fmt"

type ExceptionError struct {
	Type      string
	Value     string
	Traceback string
}

func (e ExceptionError) Error() string {
	return fmt.Sprintf(
		"python caught exception:\n%s\n%s\nTraceback:\n%s",
		e.Type,
		e.Value,
		e.Traceback,
	)
}

func LastError() error {
	var (
		pyType      *Object
		pyValue     *Object
		pyTraceback *Object
	)

	//nolint:gocritic
	C.PyErr_Fetch(&pyType, &pyValue, &pyTraceback)

	defer func() {
		pyType.DecRef()
		pyValue.DecRef()
		pyTraceback.DecRef()
	}()

	if pyType == nil {
		return nil
	}

	return ExceptionError{
		Type:      FromString(pyType),
		Value:     FromString(pyValue),
		Traceback: FromString(pyTraceback),
	}
}
