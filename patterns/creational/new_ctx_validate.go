package creational

import "context"

func NewWithContextValidate[
	O ObjectValidable,
	FO OptionWithContext[O] | optionWithContextAlias[O],
](
	ctx context.Context,
	opts ...FO,
) (O, error) {
	var obj O

	return ConstructWithContextValidate(ctx, obj, opts...)
}

func ConstructWithContextValidate[
	O ObjectValidable,
	FO OptionWithContext[O] | optionWithContextAlias[O],
](
	ctx context.Context,
	obj O,
	opts ...FO,
) (O, error) {
	obj, err := ConstructWithContext(ctx, obj, opts...)
	if err != nil {
		return obj, err
	}

	if validateErr := obj.Validate(); validateErr != nil {
		return obj, validateErr
	}

	return obj, nil
}
