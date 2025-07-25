package require

// TODO(a.gurinov): Deal with msgAndArgs vs (msg, args...)

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// TODO(a.gurinov): Deal with messages below instead of msgAndArgs...
// "actual content of YAML document doesn't match expected one; i: %d; %s", i, userMsg,
// "unexpected error while decoding YAML document; file: %s; i: %d; %s",

func YAMLEq(t *testing.T, expected string, actual string, msgAndArgs ...any) {
	t.Helper()

	expectedDocuments := decodeYAMLDocuments(t, expected, msgAndArgs...)
	actualDocuments := decodeYAMLDocuments(t, actual, msgAndArgs...)

	require.Equal(
		t,
		len(expectedDocuments),
		len(actualDocuments),
		"expected and actual number of decoded YAML documents differ",
	)

	for i := range expectedDocuments {
		require.Equal(t, yaml.DocumentNode, expectedDocuments[i])
		require.Equal(t, yaml.DocumentNode, actualDocuments[i])
		require.Equal(
			t,
			expectedDocuments[i],
			actualDocuments[i],
			msgAndArgs...,
		)
	}
}

func YAMLFilesEq(t *testing.T, expectedFilepath string, actualFilepath string, msgAndArgs ...any) {
	t.Helper()

	expectedBytes, expectedErr := os.ReadFile(expectedFilepath)
	require.NoError(t, expectedErr, msgAndArgs...)

	actualBytes, actualErr := os.ReadFile(actualFilepath)
	require.NoError(t, actualErr, msgAndArgs...)

	YAMLEq(
		t,
		string(expectedBytes),
		string(actualBytes),
		msgAndArgs...,
	)
}

func decodeYAMLDocuments(t *testing.T, content string, msgAndArgs ...any) []yaml.Node {
	t.Helper()

	var (
		doc     yaml.Node
		docs    []yaml.Node
		decoder = yaml.NewDecoder(strings.NewReader(content))
	)

	for i := 0; ; i++ {
		err := decoder.Decode(&doc)
		if errors.Is(err, io.EOF) {
			break
		}

		require.NoError(t, err, msgAndArgs...)

		docs = append(docs, doc)
	}

	return docs
}
