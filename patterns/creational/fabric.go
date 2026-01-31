package creational

import "context"

type (
	Fabric[
		O Object,
		FO Option[O] | optionAlias[O],
	] []FO

	FabricWithValidate[
		O ObjectValidable,
		FO Option[O] | optionAlias[O],
	] []FO

	FabricWithContext[
		O Object,
		FO OptionWithContext[O] | optionWithContextAlias[O],
	] []FO

	FabricWithContextValidate[
		O ObjectValidable,
		FO OptionWithContext[O] | optionWithContextAlias[O],
	] []FO
)

func (f Fabric[O, FO]) New(opts ...FO) (O, error) {
	opts = append(f, opts...)

	return New(opts...)
}

func (f FabricWithValidate[O, FO]) New(opts ...FO) (O, error) {
	opts = append(f, opts...)

	return NewWithValidate(opts...)
}

func (f FabricWithContext[O, FO]) New(ctx context.Context, opts ...FO) (O, error) {
	opts = append(f, opts...)

	return NewWithContext(ctx, opts...)
}

func (f FabricWithContextValidate[O, FO]) New(ctx context.Context, opts ...FO) (O, error) {
	opts = append(f, opts...)

	return NewWithContextValidate(ctx, opts...)
}
