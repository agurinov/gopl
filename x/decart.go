package x

func DecartDuo[A any, B any, E any](
	a []A,
	b []B,
	mergeF func(A, B) E,
) []E {
	out := make([]E, 0, len(a)*len(b))

	for i := range a {
		for ii := range b {
			out = append(out, mergeF(a[i], b[ii]))
		}
	}

	return out
}
