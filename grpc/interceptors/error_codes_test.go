package interceptors_test

import (
	"context"
	"database/sql"
	"io"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/agurinov/gopl/grpc/interceptors"
	pl_testing "github.com/agurinov/gopl/testing"
)

func TestStruct_Method(t *testing.T) {
	pl_testing.Init(t)

	type (
		args struct {
			handler grpc.UnaryHandler
		}
		results struct {
			out any
		}
		S struct {
			A string `validate:"required"`
		}
	)

	ctx := t.Context()

	cases := map[string]struct {
		args    args
		results results
		pl_testing.TestCase
	}{
		"case00: no error": {
			args: args{
				handler: func(ctx context.Context, req any) (any, error) {
					return "out1", nil
				},
			},
			results: results{
				out: "out1",
			},
		},
		"case01: sql not found": {
			args: args{
				handler: func(ctx context.Context, req any) (any, error) {
					return "out2", sql.ErrNoRows
				},
			},
			results: results{},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: status.Error(codes.NotFound, "<redacted>"),
			},
		},
		"case02: validation error": {
			args: args{
				handler: func(ctx context.Context, req any) (any, error) {
					return "out3", validator.New().Struct(S{})
				},
			},
			results: results{},
			TestCase: pl_testing.TestCase{
				MustFail: true,
				MustFailIsErr: status.Error(
					codes.InvalidArgument,
					"Key: 'S.A' Error:Field validation for 'A' failed on the 'required' tag",
				),
			},
		},
		"case03: already grpc message": {
			args: args{
				handler: func(ctx context.Context, req any) (any, error) {
					return "out4", status.Error(codes.FailedPrecondition, "foobar")
				},
			},
			results: results{},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: status.Error(codes.FailedPrecondition, "foobar"),
			},
		},
		"case04: unknown error": {
			args: args{
				handler: func(ctx context.Context, req any) (any, error) {
					return "out5", io.EOF
				},
			},
			results: results{},
			TestCase: pl_testing.TestCase{
				MustFail:      true,
				MustFailIsErr: status.Error(codes.Unknown, "<redacted>"),
			},
		},
	}

	for name := range cases {
		tc := cases[name]

		t.Run(name, func(t *testing.T) {
			tc.Init(t)

			out, err := interceptors.ErrorConverterUnaryServer(
				ctx,
				"in",
				new(grpc.UnaryServerInfo),
				tc.args.handler,
			)
			tc.CheckError(t, err)
			require.Equal(t, tc.results.out, out)
		})
	}
}
