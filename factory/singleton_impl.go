package pl_factory

import (
	"sync"
)

type singleton[O object] struct {
	factory Factory[O]
	object  O
	err     error
	once    sync.Once
}

func (s *singleton[O]) NewObject() (O, error) {
	s.once.Do(func() {
		s.object, s.err = s.factory.NewObject()
	})

	return s.object, s.err
}

func (s *singleton[O]) MustNewObject() O {
	s.once.Do(func() {
		s.object = s.factory.MustNewObject()
	})

	return s.object
}

func NewSingleton[O object](opts ...Option[O]) Factory[O] {
	f := factory[O]{
		opts: opts,
	}

	return &singleton[O]{
		factory: f,
	}
}
