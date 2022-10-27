package pl_option

type Constraint interface {
	any
}

type Option[T Constraint] func(*T) error

func New[T Constraint](opts ...Option[T]) (T, error) {
	var t T

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if err := opt(&t); err != nil {
			return t, err
		}
	}

	return t, nil
}
