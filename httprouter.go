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
	m *http.ServeMux
	// NotFoundHandler is the handler to call when the router receives a
	// request for a path that is not registered with any handler. Defaults to
	// http.NotFoundHandler.
	NotFoundHandler http.Handler
}

// New returns a new HTTP Router.
func New() *Router {
	return &Router{m: http.NewServeMux(), NotFoundHandler: http.NotFoundHandler()}
}

// Handle registers a new handler with given method and path pattern. Responds
// to the client with a 500 status code if the handler returns an error. Use
// r.PathValue() in given Handler to retrieve path parameters from the request.
//
// See https://pkg.go.dev/net/http#ServeMux for more details regarding pattern
// matching, precedence, etc.
func (r *Router) Handle(method, pattern string, h Handler, mw ...Middleware) {
	h = wrapMiddleware(mw, h)
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
