package errors

import (
	"context"
	"errors"
	"io"
)

type timeout interface {
	Timeout() bool
}

type temporary interface {
	Temporary() bool
}

func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, io.EOF) {
		return true
	}

	if errors.Is(err, context.Canceled) {
		return true
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var timeoutError timeout
	if errors.As(err, &timeoutError) {
		if timeoutError.Timeout() {
			return true
		}
	}

	var temporaryError temporary
	if errors.As(err, &temporaryError) {
		if temporaryError.Temporary() {
			return true
		}
	}

	return false
}
