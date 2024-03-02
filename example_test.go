package httprouter_test

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/zorcal/httprouter"
)

func Example() {
	r := httprouter.New()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusNotFound)
	})

	logMiddleware := func(h httprouter.Handler) httprouter.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			now := time.Now()

			err := h(w, r)

			slog.Log(r.Context(), slog.LevelInfo, "request info",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(now).String())

			return err
		}
	}

	handleIndex := func(w http.ResponseWriter, r *http.Request) error {
		fmt.Fprintf(w, "Hello World!")
		return nil
	}
	r.Handle(http.MethodGet, "/{$}", handleIndex, logMiddleware)
}
