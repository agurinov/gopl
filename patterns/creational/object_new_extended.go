package creational

import "context"

func NewExtended[
	O ObjectExtended,
	FO OptionWithContext[O] | OptionFuncWithContext[O],
](
	ctx context.Context,
	opts ...FO,
) (O, error) {
	var obj O

	return ConstructExtended(ctx, obj, opts...)
}

func ConstructExtended[
	O ObjectExtended,
	FO OptionWithContext[O] | OptionFuncWithContext[O],
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

	if initErr := obj.Init(ctx); initErr != nil {
		return obj, initErr
	}

	return obj, nil
}
