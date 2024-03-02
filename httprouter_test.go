package httprouter_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zorcal/httprouter"
)

func TestRouter(t *testing.T) {
	buf := bytes.Buffer{}

	globalMw := func(h httprouter.Handler) httprouter.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			buf.WriteString("global.")
			return h(w, r)
		}
	}
	groupMw := func(h httprouter.Handler) httprouter.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			buf.WriteString("group.")
			return h(w, r)
		}
	}
	firstMw := func(h httprouter.Handler) httprouter.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			buf.WriteString("first.")
			return h(w, r)
		}
	}
	secondMw := func(h httprouter.Handler) httprouter.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			buf.WriteString("second.")
			return h(w, r)
		}
	}

	h := func(w http.ResponseWriter, r *http.Request) error {
		buf.WriteString("handler")
		return nil
	}

	r := httprouter.New(globalMw)
	r.Handle(http.MethodGet, "/{$}", h, firstMw, secondMw)

	api := r.Group("/group", groupMw)
	api.Handle(http.MethodGet, "/{$}", h, firstMw, secondMw)

	srv := httptest.NewServer(r)

	if _, err := srv.Client().Get(srv.URL + "/"); err != nil {
		t.Fatalf("issue GET /: %v", err)
	}
	if got, want := buf.String(), "global.first.second.handler"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	buf.Reset()

	if _, err := srv.Client().Get(srv.URL + "/group"); err != nil {
		t.Fatalf("issue GET /group: %v", err)
	}
	if got, want := buf.String(), "global.group.first.second.handler"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	buf.Reset()
}

func TestNotFoundHandler(t *testing.T) {
	test := func(t *testing.T, r *httprouter.Router, wantBody string) {
		srv := httptest.NewServer(r)

		resp, err := srv.Client().Get(srv.URL + "/not-found")
		if err != nil {
			t.Fatalf("issue GET /not-found: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("got status code %d, want %d", resp.StatusCode, http.StatusNotFound)
		}

		slurp, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read response body: %v", err)
		}
		gotBody := string(slurp)

		if gotBody != wantBody {
			t.Errorf("got body %q, want %q", gotBody, wantBody)
		}
	}

	t.Run("default", func(t *testing.T) {
		r := httprouter.New()
		test(t, r, "404 page not found\n")
	})

	t.Run("custom", func(t *testing.T) {
		r := httprouter.New()
		r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "custom 404 page not found\n")
		})
		test(t, r, "custom 404 page not found\n")
	})
}

func TestNestedGroups(t *testing.T) {
	r := httprouter.New()
	group1 := r.Group("/group1")
	group11 := group1.Group("/group11")

	group11.Handle(http.MethodGet, "/destination", func(w http.ResponseWriter, r *http.Request) error {
		return nil
	})

	srv := httptest.NewServer(r)

	resp, err := srv.Client().Get(srv.URL + "/group1/group11/destination")
	if err != nil {
		t.Fatalf("issue GET /group1/group11/destination: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("got status code %d, want %d", resp.StatusCode, http.StatusOK)
	}
}
