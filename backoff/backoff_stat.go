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

	b.WriteString("retry ")
	b.WriteString(strconv.FormatUint(
		uint64(s.RetryIndex),
		base10,
	))
	b.WriteString("/")
	b.WriteString(strconv.FormatUint(
		uint64(s.MaxRetries),
		base10,
	))
	b.WriteString("; duration=")
	b.WriteString(s.Duration.String())

	return b.String()
}
