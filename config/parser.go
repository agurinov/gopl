package config

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"

	"github.com/agurinov/gopl/bitset"
)

type (
	Source      func(context.Context) ([]byte, error)
	Parser      func(context.Context, []byte, any) error
	flag        = uint8
	ParserFlags = bitset.BitSet[flag]
)

const (
	YAML uint8 = 1 << iota
	JSON
	TOML
	XML
	Silent
)

func Parse[T any](
	ctx context.Context,
	flags flag,
	sources ...Source,
) (T, error) {
	var (
		cfg    T
		parser Parser
	)

	switch bs := bitset.New(flags); {
	case bs.Has(YAML):
		parser = parseYAML
	default:
		return cfg, errors.New("unsupported parser")
	}

	for _, source := range sources {
		if source == nil {
			continue
		}

		data, err := source(ctx)

		switch {
		case err != nil:
			return cfg, err
		case len(data) == 0:
			continue
		}

		if parseErr := parser(ctx, data, &cfg); parseErr != nil {
			return cfg, parseErr
		}
	}

	v := validator.New()
	if err := v.RegisterValidation("notblank", validators.NotBlank); err != nil {
		return cfg, err
	}

	if err := v.Struct(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
