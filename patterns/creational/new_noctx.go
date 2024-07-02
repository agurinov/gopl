package creational

func New[
	O Object,
	FO Option[O] | optionAlias[O],
](
	opts ...FO,
) (O, error) {
	var obj O

	return Construct(obj, opts...)
}

func MustNew[
	O Object,
	FO Option[O] | optionAlias[O],
](
	opts ...FO,
) O {
	var obj O

	obj, err := Construct(obj, opts...)
	if err != nil {
		panic(err)
	}

	return obj
}

func Construct[
	O Object,
	FO Option[O] | optionAlias[O],
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
