// Package httprouter is a thin wrapper around http.ServeMux with support
// for middleware and a custom handler type that returns an error.
package httprouter

import (
	"fmt"
	"net/http"
)

// Handler is a function that can be registered to a route to handle HTTP
// requests.
type Handler func(w http.ResponseWriter, r *http.Request) error

// Router is a thin wrapper around http.ServeMux that allows for registering
// handlers with middleware for different HTTP methods and patterns.
type Router struct {
	// NotFoundHandler is the handler to call when the router receives a
	// request for a path that is not registered with any handler. Defaults to
	// http.NotFoundHandler.
	NotFoundHandler http.Handler

	mw []Middleware
	m  *http.ServeMux
}

// New returns a new HTTP Router. Middleware are executed in the order they are
// provided and before any handler-specific middleware.
func New(mw ...Middleware) *Router {
	return &Router{
		NotFoundHandler: http.NotFoundHandler(),
		mw:              mw,
		m:               http.NewServeMux(),
	}
}

// Handle registers a new handler with given method and path pattern.
// Middleware are exectued in the order they are provided and after any global
// middleware. Responds to the client with a 500 status code if the handler
// returns an error. Use [http.Request.PathValue] to retrieve path parameters
// from the request.
//
// Refer to [http.ServeMux] for how pattern matching, precedence, etc. works.
func (r *Router) Handle(method, pattern string, h Handler, mw ...Middleware) {
	h = wrapMiddleware(mw, h)
	h = wrapMiddleware(r.mw, h)
	r.m.HandleFunc(fmt.Sprintf("%s %s", method, pattern), func(w http.ResponseWriter, req *http.Request) {
		if err := h(w, req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

// ServeHTTP implements the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if _, pattern := r.m.Handler(req); pattern == "" {
		r.NotFoundHandler.ServeHTTP(w, req)
	} else {
		r.m.ServeHTTP(w, req)
	}
}
