package config

import (
	"context"
	"os"
)

func FromFile(filepath string, flags ...flag) Source {
	var bs ParserFlags
	for i := range flags {
		bs.Set(flags[i])
	}

	return func(context.Context) ([]byte, error) {
		data, err := os.ReadFile(filepath)
		if err != nil {
			if bs.Has(Silent) {
				return data, nil
			}

			return nil, err
		}

		return data, nil
	}
}
