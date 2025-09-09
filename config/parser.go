package config

import (
	"context"
	"errors"

	jsoniter "github.com/json-iterator/go"
	"github.com/knadh/koanf/maps"

	"github.com/agurinov/gopl/bitset"
)

type (
	Source    func(context.Context) ([]byte, error)
	Parser    func(context.Context, []byte, any) error
	Validator func(any) error
)

type (
	Flags       = uint8
	ParserFlags = bitset.BitSet[Flags]
)

const (
	YAML Flags = 1 << iota
	JSON
	TOML
	XML
	Silent
	NoValidate
)

func Parse[T any](
	ctx context.Context,
	flags Flags,
	sources ...Source,
) (
	T,
	error,
) {
	var (
		cfg T
		bs  = bitset.New(flags)
	)

	var (
		parser    Parser
		validator Validator
	)

	switch {
	case bs.Has(YAML):
		parser = parseYAML
	case bs.Has(JSON):
		parser = parseJSON
	default:
		return cfg, errors.New("unsupported parser")
	}

	switch {
	case bs.Has(NoValidate):
		validator = nil
	default:
		validator = validateStruct
	}

	dataMap := map[string]any{}

	for _, source := range sources {
		if source == nil {
			continue
		}

		sourceBytes, err := source(ctx)

		switch {
		case err != nil:
			return cfg, err
		case len(sourceBytes) == 0:
			continue
		}

		var patchMap map[string]any

		if pErr := parser(ctx, sourceBytes, &patchMap); pErr != nil {
			return cfg, pErr
		}

		maps.Merge(patchMap, dataMap)
	}

	// NOTE: not an elegant way, but yaml parser can read json input.
	dataBytes, err := jsoniter.Marshal(dataMap)
	if err != nil {
		return cfg, err
	}

	if pErr := parser(ctx, dataBytes, &cfg); pErr != nil {
		return cfg, pErr
	}

	if validator != nil {
		if vErr := validator(cfg); vErr != nil {
			return cfg, vErr
		}
	}

	return cfg, nil
}
