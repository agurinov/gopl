package strategies

import "time"

type Interface interface {
	Duration(retries uint32) time.Duration
}
