package py

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	c "github.com/agurinov/gopl/pattern/creational"
	"github.com/agurinov/gopl/py/cpy"
)

const (
	ListSeparator = string(filepath.ListSeparator)
)

var (
	pyInitOnce sync.Once
	pyInitErr  error //nolint:errname
)

type Interpreter struct {
	ensureVersion string
	home          []string
	path          []string
	ensureGIL     bool
}

// init creates new interpreter. Some functions allowed to call before Py_Initialize
// https://docs.python.org/3.8/c-api/init.html?highlight=pygilstate_ensure#before-python-initialization
func (i Interpreter) init() error {
	pyInitOnce.Do(func() {
		if len(i.path) != 0 {
			path := strings.Join(i.path, ListSeparator)
			if err := cpy.AppendPath(path); err != nil {
				pyInitErr = err

				return
			}
		}

		if err := cpy.Initialize(); err != nil {
			pyInitErr = err

			return
		}

		if i.ensureGIL && !cpy.ThreadsInitialized() {
			pyInitErr = errors.New("GIL is not initialized whereas expected")

			return
		}

		if i.ensureVersion != "" {
			if version := i.ShortVersion(); version != i.ensureVersion {
				pyInitErr = fmt.Errorf(
					"unexpected version of python interpreter initialized: %s; expected: %s",
					version,
					i.ensureVersion,
				)

				return
			}
		}

		// Release GIL from main thread
		cpy.SaveThread()
	})

	if pyInitErr != nil {
		return pyInitErr
	}

	return nil
}

func (i Interpreter) Close() error {
	return cpy.Finalize()
}

func (i Interpreter) ShortVersion() string {
	fullVersion := cpy.GetVersion()

	if i := strings.Index(fullVersion, " "); i != -1 {
		return fullVersion[:i]
	}

	return fullVersion
}

func (i Interpreter) Path() ([]string, error) {
	path, err := cpy.GetPath()
	if err != nil {
		return nil, err
	}

	return strings.Split(path, ListSeparator), nil
}

func (i Interpreter) Home() (string, error) {
	return cpy.GetPythonHome()
}

func (i Interpreter) GetModule(
	ctx context.Context,
	name string,
) (*Module, error) {
	gil := cpy.NewGIL()
	gil.Lock()
	defer gil.Unlock()

	return ImportModule(ctx, name)
}

func (i Interpreter) GetCallable(
	ctx context.Context,
	name string,
) (*Callable, error) {
	lastDotIndex := strings.LastIndex(name, ".")
	if lastDotIndex == -1 {
		return nil, fmt.Errorf("can't separate callable path")
	}

	var (
		moduleName   = name[:lastDotIndex]
		callableName = name[lastDotIndex+1:]
	)

	gil := cpy.NewGIL()
	gil.Lock()
	defer gil.Unlock()

	module, err := ImportModule(ctx, moduleName)
	if err != nil {
		return nil, err
	}

	defer module.Object.DecRef()

	callable, err := module.GetCallable(callableName)
	if err != nil {
		return nil, err
	}

	return callable, nil
}

func NewInterpreter(opts ...InterpreterOption) (Interpreter, error) {
	obj := Interpreter{
		ensureGIL: true,
	}

	obj, err := c.ConstructObject(obj, opts...)
	if err != nil {
		return obj, err
	}

	if err := obj.init(); err != nil {
		return obj, err
	}

	return obj, nil
}
