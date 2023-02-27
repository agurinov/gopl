package creational

type Object interface {
	any
}

func NewObject[O Object](opts ...Option[O]) (O, error) {
	var obj O

	return ConstructObject(obj, opts...)
}

func MustNewObject[O Object](opts ...Option[O]) O {
	var obj O

	return MustConstructObject(obj, opts...)
}

func ConstructObject[O Object](obj O, opts ...Option[O]) (O, error) {
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if err := opt(&obj); err != nil {
			return obj, err
		}
	}

	// TODO(a.gurinov): object must be with .Validate() method
	// Validate final object here

	return obj, nil
}

func MustConstructObject[O Object](obj O, opts ...Option[O]) O {
	obj, err := ConstructObject(obj, opts...)
	if err != nil {
		panic(err)
	}

	return obj
}
