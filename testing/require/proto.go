package pl_require

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ProtoEq(t *testing.T, expected proto.Message, actual proto.Message, msgAndArgs ...any) {
	t.Helper()

	expectedJSON, err := protojson.Marshal(expected)
	require.NoError(t, err, msgAndArgs...)

	actualJSON, err := protojson.Marshal(actual)
	require.NoError(t, err, msgAndArgs...)

	require.JSONEq(t, string(expectedJSON), string(actualJSON), msgAndArgs...)
}
