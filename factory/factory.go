package pl_factory

type Factory[O object] struct {
	opts []Option[O]
}

func (f Factory[O]) NewObject() (O, error) { return NewObject(f.opts...) }
func (f Factory[O]) MustNewObject() O      { return MustNewObject(f.opts...) }

func New[O object](opts ...Option[O]) Factory[O] {
	return Factory[O]{
		opts: opts,
	}
}
