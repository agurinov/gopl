package py

import (
	"context"
	"fmt"

	"github.com/agurinov/gopl/py/cpy"
)

type Module struct {
	*Object
	name    string
	version string
}

func ImportModule(ctx context.Context, moduleName string) (*Module, error) {
	done := make(chan struct{}, 1)
	defer close(done)

	go cpy.SetInterruptFromContext(ctx, done)

	m, err := cpy.ImportModule(moduleName)
	if err != nil {
		return nil, err
	}

	module := &Module{
		Object: m,
		name:   moduleName,
	}

	if versionAttr, err := cpy.ImportModuleItem(
		module.Object,
		cpy.ToString("__version__"),
	); err == nil {
		module.version = cpy.FromString(versionAttr)
	}

	return module, nil
}

func (m Module) GetCallable(name string) (*Callable, error) {
	callable, err := cpy.ImportModuleItem(
		m.Object,
		cpy.ToString(name),
	)
	if err != nil {
		return nil, err
	}

	if !cpy.IsCallable(callable) {
		return nil, fmt.Errorf("%s.%s is not callable", m.name, name)
	}

	return &Callable{Object: callable}, nil
}

func (m Module) Version() string { return m.version }
func (m Module) String() string  { return m.name }
