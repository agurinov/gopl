package vault

import "errors"

var (
	ErrDisabled            = errors.New("vault: feature disabled")
	ErrUnknownAuthMethod   = errors.New("vault: unknown auth method")
	ErrAmbiguousAuthMethod = errors.New("vault: ambiguous auth method: oneof approle or userpass")
)
