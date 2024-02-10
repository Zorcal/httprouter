package httprouter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/zorcal/httprouter"
)

// ExampleRouter demonstrates how the router can be used to handle HTTP requests.
func ExampleRouter() {
	r := httprouter.New()

	beforeMw := func(h httprouter.Handler) httprouter.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			fmt.Println("Before")
			return h(w, r)
		}
	}
	afterMw := func(h httprouter.Handler) httprouter.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			defer fmt.Println("After")
			return h(w, r)
		}
	}
	h := func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		fmt.Printf("Pong: %s\n", r.PathValue("pong"))
		return nil
	}
	r.Handle(http.MethodGet, "/ping/{pong}", h, beforeMw, afterMw)

	srv := httptest.NewServer(r)

	if _, err := srv.Client().Get(srv.URL + "/ping/example"); err != nil {
		fmt.Printf("issue GET request: %v\n", err)
		return
	}

	// Output:
	// Before
	// Pong: example
	// After
}

// ExampleRouter_error demonstrates how errors are handled by the router by default.
func ExampleRouter_error() {
	r := httprouter.New()

	errorMw := func(h httprouter.Handler) httprouter.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			if err := h(w, r); err != nil {
				fmt.Println("Error:", err)
				return err
			}
			return nil
		}
	}
	h := func(w http.ResponseWriter, r *http.Request) error {
		return fmt.Errorf("shit hit the fan!")
	}
	r.Handle(http.MethodGet, "/oh-no", h, errorMw)

	srv := httptest.NewServer(r)

	resp, err := srv.Client().Get(srv.URL + "/oh-no")
	if err != nil {
		fmt.Printf("issue GET request: %v\n", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)

	// Output:
	// Error: shit hit the fan!
	// Status: 500
}
