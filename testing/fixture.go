package testing

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func LoadFixture(t *testing.T, filepath string) string {
	t.Helper()

	require.NotEmpty(t, filepath)

	original, err := os.ReadFile(filepath)
	require.NoError(t, err)

	/*
		var compacted bytes.Buffer

		require.NoError(t,
			json.Compact(&compacted, original),
		)

		return compacted.String()
	*/

	return string(original)
}
