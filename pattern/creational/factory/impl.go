package factory

import (
	c "github.com/agurinov/gopl/pattern/creational"
)

type impl[O c.Object] struct {
	initial O
	opts    []c.Option[O]
}

func (f impl[O]) NewObject() (O, error) { return c.ConstructObject(f.initial, f.opts...) }
func (f impl[O]) MustNewObject() O      { return c.MustConstructObject(f.initial, f.opts...) }

func New[O c.Object](
	opts ...c.Option[impl[O]],
) (c.Factory[O], error) {
	return c.NewObject(opts...)
}

func MustNew[O c.Object](
	opts ...c.Option[impl[O]],
) c.Factory[O] {
	return c.MustNewObject(opts...)
}
