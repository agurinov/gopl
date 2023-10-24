package config

type Source[T any] func() (T, error)
