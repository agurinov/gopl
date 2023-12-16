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
