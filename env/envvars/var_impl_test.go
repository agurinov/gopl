//go:build test_unit

// TODO(a.gurinov): Set vars via env file for test.

package pl_envvars_test

import (
	"net"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	pl_envvars "github.com/agurinov/gopl.git/env/envvars"
	pl_testing "github.com/agurinov/gopl.git/testing"
)

func TestVariableStore_String(t *testing.T) {
	pl_testing.Init(t)

	require.NoError(t, os.Setenv("STRING_VALID", "foobar"))

	cases := map[string]struct {
		inputVar       pl_envvars.Variable[string]
		expectedStored string
		pl_testing.TestCase
	}{
		"nonexistence env var key": {
			inputVar: pl_envvars.String("STRING_NON_EXIST"),
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"success": {
			inputVar:       pl_envvars.String("STRING_VALID"),
			expectedStored: "foobar",
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			var stored string

			err := tc.inputVar.Store(&stored)

			tc.CheckError(t, err)

			require.Equal(t, tc.expectedStored, stored)
		})
	}
}

func TestVariableStore_Bool(t *testing.T) {
	pl_testing.Init(t)

	require.NoError(t, os.Setenv("BOOL_VALID", "true"))
	require.NoError(t, os.Setenv("BOOL_INVALID", "not a bool"))

	cases := map[string]struct {
		inputVar       pl_envvars.Variable[bool]
		expectedStored bool
		pl_testing.TestCase
	}{
		"nonexistence env var key": {
			inputVar: pl_envvars.Bool("BOOL_NON_EXIST"),
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"invalid value": {
			inputVar: pl_envvars.Bool("BOOL_INVALID"),
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"success": {
			inputVar:       pl_envvars.Bool("BOOL_VALID"),
			expectedStored: true,
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			var stored bool

			err := tc.inputVar.Store(&stored)

			tc.CheckError(t, err)

			require.Equal(t, tc.expectedStored, stored)
		})
	}
}

func TestVariableStore_Int(t *testing.T) {
	pl_testing.Init(t)

	require.NoError(t, os.Setenv("INT_VALID", "-100500"))
	require.NoError(t, os.Setenv("INT_INVALID", "not an int"))

	cases := map[string]struct {
		inputVar       pl_envvars.Variable[int]
		expectedStored int
		pl_testing.TestCase
	}{
		"nonexistence env var key": {
			inputVar: pl_envvars.Int("INT_NON_EXIST"),
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"invalid value": {
			inputVar: pl_envvars.Int("INT_INVALID"),
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"success": {
			inputVar:       pl_envvars.Int("INT_VALID"),
			expectedStored: -100500,
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			var stored int

			err := tc.inputVar.Store(&stored)

			tc.CheckError(t, err)

			require.Equal(t, tc.expectedStored, stored)
		})
	}
}

func TestVariableStore_Duration(t *testing.T) {
	pl_testing.Init(t)

	require.NoError(t, os.Setenv("DURATION_VALID", "1h2m30s"))
	require.NoError(t, os.Setenv("DURATION_INVALID", "not a duration"))

	cases := map[string]struct {
		inputVar       pl_envvars.Variable[time.Duration]
		expectedStored time.Duration
		pl_testing.TestCase
	}{
		"nonexistence env var key": {
			inputVar: pl_envvars.Duration("DURATION_NON_EXIST"),
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"invalid value": {
			inputVar: pl_envvars.Duration("DURATION_INVALID"),
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"success": {
			inputVar:       pl_envvars.Duration("DURATION_VALID"),
			expectedStored: time.Hour + 2*time.Minute + 30*time.Second,
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			var stored time.Duration

			err := tc.inputVar.Store(&stored)

			tc.CheckError(t, err)

			require.Equal(t, tc.expectedStored, stored)
		})
	}
}

func TestVariableStore_UUID(t *testing.T) {
	pl_testing.Init(t)

	require.NoError(t, os.Setenv("UUID_VALID", "711835ae-613e-4cb9-9215-f3c6a8688676"))
	require.NoError(t, os.Setenv("UUID_INVALID", "not a bool"))

	cases := map[string]struct {
		inputVar       pl_envvars.Variable[uuid.UUID]
		expectedStored uuid.UUID
		pl_testing.TestCase
	}{
		"nonexistence env var key": {
			inputVar: pl_envvars.UUID("UUID_NON_EXIST"),
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"invalid value": {
			inputVar: pl_envvars.UUID("UUID_INVALID"),
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"success": {
			inputVar:       pl_envvars.UUID("UUID_VALID"),
			expectedStored: uuid.MustParse("711835ae-613e-4cb9-9215-f3c6a8688676"),
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			var stored uuid.UUID

			err := tc.inputVar.Store(&stored)

			tc.CheckError(t, err)

			require.Equal(t, tc.expectedStored, stored)
		})
	}
}

func TestVariableStore_IP(t *testing.T) {
	pl_testing.Init(t)

	require.NoError(t, os.Setenv("IP_V4_VALID", "192.168.0.1"))
	require.NoError(t, os.Setenv("IP_V6_VALID", "2001:db8::68"))
	require.NoError(t, os.Setenv("IP_INVALID", "192.0.2"))

	cases := map[string]struct {
		inputVar       pl_envvars.Variable[net.IP]
		expectedStored net.IP
		pl_testing.TestCase
	}{
		"nonexistence env var key": {
			inputVar: pl_envvars.IP("IP_NON_EXIST"),
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"invalid value": {
			inputVar: pl_envvars.IP("IP_INVALID"),
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"success ipv4": {
			inputVar:       pl_envvars.IP("IP_V4_VALID"),
			expectedStored: net.IPv4(192, 168, 0, 1),
		},
		"success ipv6": {
			inputVar:       pl_envvars.IP("IP_V6_VALID"),
			expectedStored: net.IP{0x20, 0x1, 0xd, 0xb8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x68},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			var stored net.IP

			err := tc.inputVar.Store(&stored)

			tc.CheckError(t, err)

			require.Equal(t, tc.expectedStored, stored)
		})
	}
}

func TestVariableStore_URL(t *testing.T) {
	pl_testing.Init(t)

	require.NoError(t, os.Setenv("URL_VALID", "http://domain.com/foo/bar"))
	require.NoError(t, os.Setenv("URL_INVALID", "://$? this is not an URL ://"))

	cases := map[string]struct {
		inputVar       pl_envvars.Variable[url.URL]
		expectedStored url.URL
		pl_testing.TestCase
	}{
		"nonexistence env var key": {
			inputVar: pl_envvars.URL("URL_NON_EXIST"),
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"invalid value": {
			inputVar: pl_envvars.URL("URL_INVALID"),
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"success": {
			inputVar: pl_envvars.URL("URL_VALID"),
			expectedStored: url.URL{
				Scheme: "http",
				Host:   "domain.com",
				Path:   "/foo/bar",
			},
		},
	}

	for name, tc := range cases {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			var stored url.URL

			err := tc.inputVar.Store(&stored)

			tc.CheckError(t, err)

			require.Equal(t, tc.expectedStored, stored)
		})
	}
}
