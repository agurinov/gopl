package creational

import "context"

func NewWithContext[
	O Object,
	FO OptionWithContext[O] | optionWithContextAlias[O],
](
	ctx context.Context,
	opts ...FO,
) (O, error) {
	var obj O

	return ConstructWithContext(ctx, obj, opts...)
}

func ConstructWithContext[
	O Object,
	FO OptionWithContext[O] | optionWithContextAlias[O],
](
	ctx context.Context,
	obj O,
	opts ...FO,
) (O, error) {
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if err := opt(ctx, &obj); err != nil {
			return obj, err
		}
	}

	return obj, nil
}
