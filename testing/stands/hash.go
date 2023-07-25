package stands

import (
	"crypto"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	maxHashLen = 8
)

func hash(t *testing.T) string {
	t.Helper()

	hasher := crypto.MD5.New()
	require.NotNil(t, hasher)

	_, _ = hasher.Write([]byte(t.Name()))
	hashhex := hex.EncodeToString(hasher.Sum(nil))

	return hashhex[:maxHashLen]
}
