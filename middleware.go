package httprouter

// Middleware is a handler function designed to run code before and/or after
// another Handler.
type Middleware func(Handler) Handler

func wrapMiddleware(mw []Middleware, h Handler) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		wrap := mw[i]
		if wrap != nil {
			h = wrap(h)
		}
	}
	return h
}
