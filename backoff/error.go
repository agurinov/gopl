package backoff

import (
	"strconv"
	"strings"
)

const (
	base10 = 10
)

type RetryLimitError struct {
	BackoffName string
	MaxRetries  uint32
}

func (e RetryLimitError) Error() string {
	var b strings.Builder

	b.WriteString("backoff(")
	b.WriteString(e.BackoffName)
	b.WriteString(") retry limit exceeded; max_retries=")
	b.WriteString(strconv.FormatUint(
		uint64(e.MaxRetries),
		base10,
	))

	return b.String()
}
