package kafka

import "errors"

var (
	ErrEmptyConsumerLibrary = errors.New("<nil> consumer library")
	ErrEmptySerializer      = errors.New("<nil> event serializer")
	ErrInvalidConfig        = errors.New("invalid config")
	ErrInvalidConfigMap     = errors.New("invalid configmap")
	ErrInvalidEventHandler  = errors.New("invalid event handler")
	ErrInvalidEventPosition = errors.New("invalid event position")
	ErrUnknownSASLMechanism = errors.New("unknown SASL mechanism")
)
