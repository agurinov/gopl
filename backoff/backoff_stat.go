package backoff

import (
	"strconv"
	"strings"
	"time"
)

type Stat struct {
	Duration   time.Duration
	RetryIndex uint32
	MaxRetries uint32
}

var EmptyStat = Stat{}

func (s Stat) String() string {
	var b strings.Builder

	_, _ = b.WriteString("retry ")
	_, _ = b.WriteString(strconv.FormatUint(
		uint64(s.RetryIndex),
		base10,
	))
	_, _ = b.WriteString("/")
	_, _ = b.WriteString(strconv.FormatUint(
		uint64(s.MaxRetries),
		base10,
	))
	_, _ = b.WriteString("; duration=")
	_, _ = b.WriteString(s.Duration.String())

	return b.String()
}
