package pl_option

func WithA[T Constraint](a int) Option[T] {
	return func(t *T) error {
		// t.a = a

		return nil
	}
}
