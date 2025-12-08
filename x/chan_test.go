package x_test

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	pl_testing "github.com/agurinov/gopl/testing"
	"github.com/agurinov/gopl/x"
)

func TestFlattenErrors(t *testing.T) {
	pl_testing.Init(t)

	t.Run("ch1", func(t *testing.T) {
		pl_testing.Init(t)

		ch := make(chan error, 10)

		t.Cleanup(func() { close(ch) })

		ch <- io.EOF
		ch <- io.ErrNoProgress

		err := x.FlattenErrors(ch)

		require.ErrorIs(t, err, io.EOF)
		require.ErrorIs(t, err, io.ErrNoProgress)

		asList := errorAsList(err)
		require.Len(t, asList, 2)
	})

	t.Run("ch2: empty", func(t *testing.T) {
		pl_testing.Init(t)

		ch := make(chan error, 10)

		t.Cleanup(func() { close(ch) })

		err := x.FlattenErrors(ch)

		require.NoError(t, err)
	})

	t.Run("ch2: closed", func(t *testing.T) {
		pl_testing.Init(t)

		ch := make(chan error, 10)
		ch <- io.EOF
		close(ch)

		err := x.FlattenErrors(ch)

		require.Error(t, err)

		asList := errorAsList(err)
		require.Len(t, asList, 1)
	})

	t.Run("combined", func(t *testing.T) {
		pl_testing.Init(t)

		ch1 := make(chan error, 10)
		ch2 := make(chan error, 10)
		ch3 := make(chan error, 10)

		t.Cleanup(func() { close(ch1) })
		t.Cleanup(func() { close(ch2) })
		t.Cleanup(func() { close(ch3) })

		ch1 <- io.EOF
		ch1 <- io.ErrNoProgress

		ch3 <- errors.New("3")
		ch3 <- errors.New("4")
		ch3 <- errors.New("5")

		err := x.FlattenErrors(ch1, ch2, ch3)

		require.Error(t, err)

		asList := errorAsList(err)
		require.Len(t, asList, 5)
	})
}

func errorAsList(err error) []error {
	typed, ok := err.(interface{ Unwrap() []error })
	if !ok {
		return nil
	}

	return typed.Unwrap()
}
