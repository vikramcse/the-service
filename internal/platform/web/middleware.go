package web

// Middleware is a function designed to run some code before and/or after
// another Handler. It is designed to remove boilerplate or other conserns
// not direct to any given Handler
type Middleware func(Handler) Handler

// wrapMiddleare creates a new handler by wrapping middleware around a final
// handler. The middlware handlers will be executed by requests in the order
// they are privided
func wrapMiddleware(mv []Middleware, handler Handler) Handler {
	for i := len(mv) - 1; i >= 0; i-- {
		h := mv[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
