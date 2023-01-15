package factory

import (
	c "github.com/agurinov/gopl.git/pattern/creational"
)

func WithInitialObject[O c.Object](initial O) c.Option[impl[O]] {
	return func(f *impl[O]) error {
		f.initial = initial

		return nil
	}
}

func WithOptions[O c.Object](opts ...c.Option[O]) c.Option[impl[O]] {
	return func(f *impl[O]) error {
		f.opts = opts

		return nil
	}
}
