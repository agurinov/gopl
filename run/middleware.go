package run

type (
	Middleware[H any]  = func(H) H
	Middlewares[H any] []Middleware[H]
)

func (mws Middlewares[H]) Handler(handler H) H {
	if len(mws) == 0 {
		return handler
	}

	for i := len(mws) - 1; i >= 0; i-- {
		handler = mws[i](handler)
	}

	return handler
}
