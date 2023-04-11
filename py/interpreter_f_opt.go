package py

import (
	"path/filepath"

	c "github.com/agurinov/gopl/pattern/creational"
)

type InterpreterOption = c.Option[Interpreter]

func WithPythonPath(pyPath ...string) InterpreterOption {
	return func(i *Interpreter) error {
		for i := range pyPath {
			absPath, err := filepath.Abs(pyPath[i])
			if err != nil {
				return err
			}

			pyPath[i] = absPath
		}

		i.path = pyPath

		return nil
	}
}

func WithPythonHome(pyHome ...string) InterpreterOption {
	return func(i *Interpreter) error {
		i.home = pyHome

		return nil
	}
}

func WithEnsureGIL(g bool) InterpreterOption {
	return func(i *Interpreter) error {
		i.ensureGIL = g

		return nil
	}
}

func WithEnsureVersion(v string) InterpreterOption {
	return func(i *Interpreter) error {
		i.ensureVersion = v

		return nil
	}
}
