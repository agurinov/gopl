package pl_envvars

import (
	"strconv"
	"time"
)

var (
	toStringMapper   = func(val string) (string, error) { return val, nil }
	toBoolMapper     = strconv.ParseBool
	toIntMapper      = strconv.Atoi
	toDurationMapper = time.ParseDuration
)
