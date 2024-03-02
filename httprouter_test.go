package httprouter_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zorcal/httprouter"
)

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
