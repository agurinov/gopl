package creational

func NewWithValidate[
	O ObjectValidable,
	FO Option[O] | optionAlias[O],
](
	opts ...FO,
) (O, error) {
	var obj O

	return ConstructWithValidate(obj, opts...)
}

func ConstructWithValidate[
	O ObjectValidable,
	FO Option[O] | optionAlias[O],
](
	obj O,
	opts ...FO,
) (O, error) {
	obj, err := Construct(obj, opts...)
	if err != nil {
		return obj, err
	}

	if validateErr := obj.Validate(); validateErr != nil {
		return obj, validateErr
	}

	return obj, nil
}
