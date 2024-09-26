package config

import (
	"bytes"
	"context"

	jsoniter "github.com/json-iterator/go"
)

func parseJSON(
	_ context.Context,
	data []byte,
	obj any,
) error {
	var (
		r       = bytes.NewBuffer(data)
		decoder = jsoniter.NewDecoder(r)
	)

	decoder.DisallowUnknownFields()

	return decoder.Decode(obj)
}
