//go:build test_unit

package pl_require_test

import (
	"testing"

	pl_testing "github.com/agurinov/gopl.git/testing"
	pl_require "github.com/agurinov/gopl.git/testing/require"
)

func TestYAMLFilesEq(t *testing.T) {
	pl_testing.Init(t)

	cases := map[string]struct {
		oneFilepath string
		twoFilepath string
		pl_testing.TestCase
	}{
		"different files": {
			oneFilepath: "testdata/2doc.yaml",
			twoFilepath: "testdata/3doc.yaml",
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"equal": {
			oneFilepath: "testdata/2doc.yaml",
			twoFilepath: "testdata/2doc.yaml",
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			pl_require.YAMLFilesEq(t, tc.oneFilepath, tc.twoFilepath)
		})
	}
}
