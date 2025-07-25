package envvars_test

import (
	"net"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"

	envvars "github.com/agurinov/gopl/env/envvars"
	pl_testing "github.com/agurinov/gopl/testing"
)

type (
	args[T envvars.T] struct {
		envvar envvars.Variable[T]
	}
	results[T envvars.T] struct {
		stored T
	}
	tc[T envvars.T] struct {
		args    args[T]
		results results[T]
		pl_testing.TestCase
	}
)

func RunCases[T envvars.T](
	t *testing.T,
	cases map[string]tc[T],
) {
	t.Helper()

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			var (
				stored T
				err    = tc.args.envvar.Store(&stored)
			)

			tc.CheckError(t, err)
			require.Equal(t, tc.results.stored, stored)
		})
	}
}

func TestVariableStore(t *testing.T) {
	pl_testing.Init(t)

	require.NoError(t, godotenv.Load("testdata/.env"))

	t.Run("string", func(t *testing.T) {
		cases := map[string]tc[string]{
			"case00: nonexistence env var key": {
				args: args[string]{
					envvar: envvars.String("STRING_NON_EXIST"),
				},
				TestCase: pl_testing.TestCase{
					MustFail:      true,
					MustFailIsErr: envvars.ErrNoVar,
				},
			},
			"case01: success": {
				args: args[string]{
					envvar: envvars.String("STRING_VALID"),
				},
				results: results[string]{
					stored: "foobar",
				},
			},
		}

		RunCases(t, cases)
	})

	t.Run("bool", func(t *testing.T) {
		cases := map[string]tc[bool]{
			"case00: nonexistence env var key": {
				args: args[bool]{
					envvar: envvars.Bool("BOOL_NON_EXIST"),
				},
				TestCase: pl_testing.TestCase{
					MustFail:      true,
					MustFailIsErr: envvars.ErrNoVar,
				},
			},
			"case01: invalid value": {
				args: args[bool]{
					envvar: envvars.Bool("BOOL_INVALID"),
				},
				TestCase: pl_testing.TestCase{
					MustFail:      true,
					MustFailIsErr: strconv.ErrSyntax,
				},
			},
			"case02: success": {
				args: args[bool]{
					envvar: envvars.Bool("BOOL_VALID"),
				},
				results: results[bool]{
					stored: true,
				},
			},
		}

		RunCases(t, cases)
	})

	t.Run("int", func(t *testing.T) {
		cases := map[string]tc[int]{
			"case00: nonexistence env var key": {
				args: args[int]{
					envvar: envvars.Int("INT_NON_EXIST"),
				},
				TestCase: pl_testing.TestCase{
					MustFail:      true,
					MustFailIsErr: envvars.ErrNoVar,
				},
			},
			"case01: invalid value": {
				args: args[int]{
					envvar: envvars.Int("INT_INVALID"),
				},
				TestCase: pl_testing.TestCase{
					MustFail:      true,
					MustFailIsErr: strconv.ErrSyntax,
				},
			},
			"case02: success": {
				args: args[int]{
					envvar: envvars.Int("INT_VALID"),
				},
				results: results[int]{
					stored: 100500,
				},
			},
		}

		RunCases(t, cases)
	})

	t.Run("duration", func(t *testing.T) {
		cases := map[string]tc[time.Duration]{
			"case00: nonexistence env var key": {
				args: args[time.Duration]{
					envvar: envvars.Duration("DURATION_NON_EXIST"),
				},
				TestCase: pl_testing.TestCase{
					MustFail:      true,
					MustFailIsErr: envvars.ErrNoVar,
				},
			},
			"case01: invalid value": {
				args: args[time.Duration]{
					envvar: envvars.Duration("DURATION_INVALID"),
				},
				TestCase: pl_testing.TestCase{
					MustFail: true,
				},
			},
			"case02: success": {
				args: args[time.Duration]{
					envvar: envvars.Duration("DURATION_VALID"),
				},
				results: results[time.Duration]{
					stored: time.Hour + 2*time.Minute + 30*time.Second,
				},
			},
		}

		RunCases(t, cases)
	})

	t.Run("uuid", func(t *testing.T) {
		cases := map[string]tc[uuid.UUID]{
			"case00: nonexistence env var key": {
				args: args[uuid.UUID]{
					envvar: envvars.UUID("UUID_NON_EXIST"),
				},
				TestCase: pl_testing.TestCase{
					MustFail:      true,
					MustFailIsErr: envvars.ErrNoVar,
				},
			},
			"case01: invalid value": {
				args: args[uuid.UUID]{
					envvar: envvars.UUID("UUID_INVALID"),
				},
				TestCase: pl_testing.TestCase{
					MustFail: true,
				},
			},
			"case02: success": {
				args: args[uuid.UUID]{
					envvar: envvars.UUID("UUID_VALID"),
				},
				results: results[uuid.UUID]{
					stored: uuid.MustParse("711835ae-613e-4cb9-9215-f3c6a8688676"),
				},
			},
		}

		RunCases(t, cases)
	})

	t.Run("ip", func(t *testing.T) {
		cases := map[string]tc[net.IP]{
			"case00: nonexistence env var key": {
				args: args[net.IP]{
					envvar: envvars.IP("IP_NON_EXIST"),
				},
				TestCase: pl_testing.TestCase{
					MustFail:      true,
					MustFailIsErr: envvars.ErrNoVar,
				},
			},
			"case01: invalid value": {
				args: args[net.IP]{
					envvar: envvars.IP("IP_INVALID"),
				},
				TestCase: pl_testing.TestCase{
					MustFail:      true,
					MustFailIsErr: envvars.ErrParseIP,
				},
			},
			"case02: success ipv4": {
				args: args[net.IP]{
					envvar: envvars.IP("IP_V4_VALID"),
				},
				results: results[net.IP]{
					stored: net.IPv4(192, 168, 0, 1),
				},
			},
			"case03: success ipv6": {
				args: args[net.IP]{
					envvar: envvars.IP("IP_V6_VALID"),
				},
				results: results[net.IP]{
					stored: net.IP{
						0x20, 0x1, 0xd, 0xb8, 0x0, 0x0, 0x0, 0x0,
						0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x68,
					},
				},
			},
		}

		RunCases(t, cases)
	})

	t.Run("url", func(t *testing.T) {
		urlParseError := new(url.Error)

		cases := map[string]tc[url.URL]{
			"case00: nonexistence env var key": {
				args: args[url.URL]{
					envvar: envvars.URL("URL_NON_EXIST"),
				},
				TestCase: pl_testing.TestCase{
					MustFail:      true,
					MustFailIsErr: envvars.ErrNoVar,
				},
			},
			"case01: invalid value": {
				args: args[url.URL]{
					envvar: envvars.URL("URL_INVALID"),
				},
				TestCase: pl_testing.TestCase{
					MustFail:      true,
					MustFailAsErr: &urlParseError,
				},
			},
			"case02: success": {
				args: args[url.URL]{
					envvar: envvars.URL("URL_VALID"),
				},
				results: results[url.URL]{
					stored: url.URL{
						Scheme: "http",
						Host:   "domain.com",
						Path:   "/foo/bar",
					},
				},
			},
		}

		RunCases(t, cases)
	})
}
