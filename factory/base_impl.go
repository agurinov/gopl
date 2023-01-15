package pl_factory

type factory[O object] struct {
	opts []Option[O]
}

func (f factory[O]) NewObject() (O, error) { return NewObject(f.opts...) }
func (f factory[O]) MustNewObject() O      { return MustNewObject(f.opts...) }

func New[O object](opts ...Option[O]) Factory[O] {
	return factory[O]{
		opts: opts,
	}
}
