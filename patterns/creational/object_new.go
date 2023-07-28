package creational

func New[
	O Object,
	FO Option[O] | OptionFunc[O],
](
	opts ...FO,
) (O, error) {
	var obj O

	return Construct(obj, opts...)
}

func MustNew[
	O Object,
	FO Option[O] | OptionFunc[O],
](
	opts ...FO,
) O {
	var obj O

	return MustConstruct(obj, opts...)
}

func Construct[
	O Object,
	FO Option[O] | OptionFunc[O],
](
	obj O,
	opts ...FO,
) (O, error) {
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

func MustConstruct[
	O Object,
	FO Option[O] | OptionFunc[O],
](
	obj O,
	opts ...FO,
) O {
	obj, err := Construct(obj, opts...)
	if err != nil {
		panic(err)
	}

	return obj
}
