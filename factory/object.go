package pl_factory

type object interface {
	any
}

type (
	ObjectConstructor[O object]     func(...Option[O]) (O, error)
	ObjectConstructorMust[O object] func(...Option[O]) O
)

func NewObject[O object](opts ...Option[O]) (O, error) {
	var obj O

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if err := opt(&obj); err != nil {
			return obj, err
		}
	}

	return obj, nil
}

func MustNewObject[O object](opts ...Option[O]) O {
	obj, err := NewObject(opts...)
	if err != nil {
		panic(err)
	}

	return obj
}
