package config

import (
	"bytes"
	"context"

	"gopkg.in/yaml.v3"
)

func parseYAML(
	_ context.Context,
	data []byte,
	obj any,
) error {
	var (
		r       = bytes.NewBuffer(data)
		decoder = yaml.NewDecoder(r)
	)

	decoder.KnownFields(true)

	return decoder.Decode(obj)
}
