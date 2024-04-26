package config

import (
	"context"

	"github.com/goccy/go-yaml"
)

func parseYAML(
	ctx context.Context,
	data []byte,
	obj any,
) error {
	return yaml.UnmarshalContext(
		ctx, data, obj,
		yaml.DisallowDuplicateKey(),
		yaml.DisallowUnknownField(),
	)
}
