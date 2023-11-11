package envvars

import "errors"

var (
	ErrNoVar    = errors.New("envvar doesn't present")
	ErrNoMapper = errors.New("envvar doesn't have mapper")
	ErrParseIP  = errors.New("can't parse IP")
)
