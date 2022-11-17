package pl_envvars

import (
	"io"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var (
	toStringMapper   = func(s string) (string, error) { return s, nil }
	toBoolMapper     = strconv.ParseBool
	toIntMapper      = strconv.Atoi
	toDurationMapper = time.ParseDuration
	toUUIDMapper     = uuid.Parse
	toIPMapper       = func(s string) (net.IP, error) {
		ip := net.ParseIP(s)
		if ip == nil {
			// TODO(a.gurinov): fix this error
			return nil, io.EOF
		}

		return ip, nil
	}
	toURLMapper = func(s string) (url.URL, error) {
		urlPtr, err := url.Parse(s)
		if err != nil {
			return url.URL{}, err
		}

		return *urlPtr, nil
	}
)
