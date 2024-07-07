package telegram_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/agurinov/gopl/telegram"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestAuth_AuthFunc(t *testing.T) {
	type (
		args struct {
			initDataString string
			dummyEnabled   bool
		}
		results struct {
			user telegram.User
		}
	)

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: no header": {
			args: args{
				initDataString: "",
			},
			results: results{},
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"case01: wrong token": {
			args: args{
				initDataString: "tma foobar",
			},
			results: results{},
			TestCase: pl_testing.TestCase{
				MustFail: true,
			},
		},
		"case02: dummy: no header": {
			args: args{
				initDataString: "",
				dummyEnabled:   true,
			},
			results: results{
				user: telegram.Dummy(),
			},
		},
		"case03: dummy: wrong token": {
			args: args{
				initDataString: "tma foobar",
				dummyEnabled:   true,
			},
			results: results{
				user: telegram.Dummy(),
			},
		},
	}

	for name := range cases {
		name, tc := name, cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			auth, err := telegram.NewAuth(
				telegram.WithAuthLogger(zaptest.NewLogger(t)),
				telegram.WithAuthDummy(tc.args.dummyEnabled),
				telegram.WithAuthBotTokens(map[string]string{"FooBot": "foo"}),
			)
			require.NoError(t, err)
			require.NotNil(t, auth)

			user, err := auth.AuthFunc(tc.args.initDataString)
			tc.CheckError(t, err)
			require.Equal(t, tc.results.user, user)
		})
	}
}
