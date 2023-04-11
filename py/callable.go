package py

import (
	"context"

	"github.com/agurinov/gopl/py/cpy"
)

type (
	Callable         struct{ *Object }
	MapperToPython   = func() ([]*Object, error)
	MapperFromPython = func(*Object) error
)

func (c *Callable) Call(
	ctx context.Context,
	mapperToPython MapperToPython,
	mapperFromPython MapperFromPython,
) error {
	gil := cpy.NewGIL()
	gil.Lock()
	defer gil.Unlock()

	done := make(chan struct{}, 1)
	defer close(done)

	go cpy.SetInterruptFromContext(ctx, done)

	var args []*Object

	if mapperToPython != nil {
		var mapperErr error

		args, mapperErr = mapperToPython()
		if mapperErr != nil {
			return mapperErr
		}
	}

	argsTuple, err := cpy.Tuple(args...)
	if err != nil {
		return err
	}

	defer argsTuple.DecRef()

	res, err := cpy.CallObject(
		c.Object,
		argsTuple,
	)
	if err != nil {
		return err
	}

	defer res.DecRef()

	if mapperFromPython != nil {
		if err := mapperFromPython(res); err != nil {
			return err
		}
	}

	return nil
}
