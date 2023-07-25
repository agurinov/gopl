package backoff

import (
	"strconv"
	"strings"
)

type RetryLimitError struct {
	BackoffName string
	MaxRetries  uint32
}

const (
	base10 = 10
)

func (e RetryLimitError) Error() string {
	var b strings.Builder

	_, _ = b.WriteString("backoff(")
	_, _ = b.WriteString(e.BackoffName)
	_, _ = b.WriteString(") retry limit exceeded; max_retries=")
	_, _ = b.WriteString(strconv.FormatUint(
		uint64(e.MaxRetries),
		base10,
	))

	return b.String()
}
