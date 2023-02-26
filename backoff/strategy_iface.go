package backoff

import "time"

type Strategy interface {
	Duration(retries uint32) time.Duration
}
